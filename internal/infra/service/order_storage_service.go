package service

import (
	"errors"
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

	res, err := handler.Handle(blockHash, hookAddress, ordersSlot)
	if err != nil {
		return nil, err
	}

	arrayLength := new(big.Int).SetBytes(common.FromHex(res.Response))
	if arrayLength.Sign() == 0 {
		return nil, ErrNoOrdersFound
	}

	orders := make([]*domain.Order, 0, arrayLength.Int64())
	slotHash := crypto.Keccak256Hash(ordersSlot.Bytes())

	for i := int64(0); i < arrayLength.Int64(); i++ {
		var orderRawData [3]uint256.Int

		for j := 0; j < 3; j++ {
			data, err := handler.Handle(blockHash, hookAddress, slotHash)
			if err != nil {
				return nil, err
			}

			orderRawData[j] = *uint256.MustFromBig(new(big.Int).SetBytes(common.FromHex(data.Response)))
			slotHash = common.BigToHash(new(big.Int).Add(new(big.Int).SetBytes(slotHash.Bytes()), big.NewInt(1)))
		}

		status, err := s.FindOrderStatus(hookAddress, big.NewInt(i), blockHash, statusSlot)
		if err != nil {
			return nil, err
		}

		orderStatus := domain.OrderStatusOpen
		if *status {
			orderStatus = domain.OrderStatusFulFilledOrClosed
		}

		order, err := domain.NewOrder(
			uint64(i+1), // the index inside of the dApp is 1-based index, instead of 0-based index from the blockchain
			hookAddress,
			&orderRawData[1],
			&orderRawData[2],
			nil,
			&orderStatus,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrderStorageService) FindOrderStatus(hookAddress common.Address, orderId *big.Int, blockHash, slot common.Hash) (*bool, error) {
	handler, err := s.GioHandlerFactory.NewGioHandler(0x27)
	if err != nil {
		return nil, err
	}

	slotHash := crypto.Keccak256Hash(common.LeftPadBytes(orderId.Bytes(), 32), slot.Bytes())
	res, err := handler.Handle(blockHash, hookAddress, slotHash)
	if err != nil {
		return nil, err
	}

	status := new(big.Int).SetBytes(common.FromHex(res.Response)).Cmp(big.NewInt(1)) == 0
	return &status, nil
}
