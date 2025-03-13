package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/henriquemarlon/swapx/pkg/gio"
	"github.com/holiman/uint256"
)

var ErrNoOrdersFound = errors.New("no orders found")

type OrderStorageService struct {
	GioHandlerFactory gio.GioHandlerFactory
}

type OrderStorageServiceInterface interface {
	FindOrderStatus(hookAddress common.Address, orderId *big.Int, blockHash, slot common.Hash) (*bool, error)
	FindOrdersBySlot(hookAddress common.Address, blockHash, ordersSlot, statusSlot common.Hash) ([]*domain.Order, error)
}

func NewOrderStorageService(gioHandlerFactory gio.GioHandlerFactory) *OrderStorageService {
	return &OrderStorageService{GioHandlerFactory: gioHandlerFactory}
}

func (s *OrderStorageService) FindOrdersBySlot(hookAddress common.Address, blockHash, ordersSlot, statusSlot common.Hash) ([]*domain.Order, error) {
	handler, err := s.GioHandlerFactory.NewGioHandler(0x27)
	if err != nil {
		return nil, err
	}

	slog.Info("/====================== Looking for orders at", "slot =====================", fmt.Sprintf("> %v", new(big.Int).SetBytes(ordersSlot.Bytes())))

	res, err := handler.Handle(blockHash, hookAddress, ordersSlot)
	if err != nil {
		return nil, err
	}

	arrayLength := new(big.Int).SetBytes(common.FromHex(res.Response))
	if arrayLength.Sign() == 0 {
		return nil, ErrNoOrdersFound
	}

	slog.Info("Total orders found in storage", "count", arrayLength.Int64())

	orders := make([]*domain.Order, 0, arrayLength.Int64())
	slotHash := crypto.Keccak256Hash(ordersSlot.Bytes())

	for i := int64(0); i < arrayLength.Int64(); i++ {
		var orderRawData [4]uint256.Int

		for j := 0; j < 4; j++ {
			data, err := handler.Handle(blockHash, hookAddress, slotHash)
			if err != nil {
				return nil, err
			}

			orderRawData[j] = *uint256.MustFromBig(new(big.Int).SetBytes(common.FromHex(data.Response)))
			slotHash = common.BigToHash(new(big.Int).Add(new(big.Int).SetBytes(slotHash.Bytes()), big.NewInt(1)))
		}

		orderRawDataBytes, err := json.Marshal(orderRawData)
		if err != nil {
			return nil, err
		}
		slog.Info("Raw order data collected from base layer with", " id", i, "data", string(orderRawDataBytes))

		isCancelled, err := s.FindOrderStatus(hookAddress, big.NewInt(i), blockHash, statusSlot)
		if err != nil {
			return nil, err
		}

		isFulfilled := new(uint256.Int).Sub(&orderRawData[2], &orderRawData[3]).IsZero()

		orderStatus := domain.OrderNotCancelledOrFulfilled
		if *isCancelled || isFulfilled {
			orderStatus = domain.OrderCancelledOrFulfilled
		}

		order, err := domain.NewOrder(
			uint64(i+1), // the index inside of the dApp is 1-based index, instead of 0-based index from the blockchain
			hookAddress,
			&orderRawData[1],
			&orderRawData[2],
			&orderRawData[3],
			nil,
			&orderStatus,
		)
		if err != nil {
			return nil, err
		}

		orderBytes, err := json.Marshal(order)
		if err != nil {
			return nil, err
		}
		slog.Info("Order found", "info", string(orderBytes))

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrderStorageService) FindOrderStatus(hookAddress common.Address, orderId *big.Int, blockHash, slot common.Hash) (*bool, error) {
	handler, err := s.GioHandlerFactory.NewGioHandler(0x27)
	if err != nil {
		return nil, err
	}

	slog.Info("/====================== Looking for order status at", "slot =====================", fmt.Sprintf("> %v", new(big.Int).SetBytes(slot.Bytes())))

	slotHash := crypto.Keccak256Hash(common.BigToHash(orderId).Bytes(), slot.Bytes())
	res, err := handler.Handle(blockHash, hookAddress, slotHash)
	if err != nil {
		return nil, err
	}

	status := new(big.Int).SetBytes(common.FromHex(res.Response)).Cmp(big.NewInt(1)) == 0
	return &status, nil
}
