package service

import (
	"errors"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/henriquemarlon/swapx/pkg/gio"
	"github.com/holiman/uint256"
)

var (
	ErrNoOrdersFound = errors.New("no orders found")
)

var ROLLUP_HTTP_SERVER_URL = os.Getenv("ROLLUP_HTTP_SERVER_URL")

type HookStorageService struct{}

type HookStorageServiceInterface interface {
	FindAllOrdersBySlot(hookAddress common.Address, blockHash, slot common.Hash) ([]*domain.Order, error)
}

func NewHookStorageService() *HookStorageService {
	return &HookStorageService{}
}

func (s *HookStorageService) FindAllOrdersBySlot(hookAddress common.Address, blockHash, slot common.Hash) ([]*domain.Order, error) {
	factory := gio.NewGioHandlerFactory(ROLLUP_HTTP_SERVER_URL)
	handler, err := factory.NewGioHandler(0x27)
	if err != nil {
		return nil, err
	}

	res, err := handler.Handle(blockHash, hookAddress, slot)
	if err != nil {
		return nil, err
	}

	arrayLength := new(big.Int).SetBytes(common.FromHex(res.Response))
	if arrayLength.Sign() == 0 {
		return nil, ErrNoOrdersFound
	}

	orders := make([]*domain.Order, 0, arrayLength.Int64())
	slotHash := crypto.Keccak256Hash(slot.Bytes())

	for i := int64(0); i < arrayLength.Int64(); i++ {
		var orderRawData []string

		for j := 0; j < 3; j++ {
			data, err := handler.Handle(blockHash, hookAddress, slotHash)
			if err != nil {
				return nil, err
			}
			orderRawData = append(orderRawData, data.Response)
			slotHash = common.BigToHash(new(big.Int).Add(new(big.Int).SetBytes(slotHash.Bytes()), big.NewInt(1)))
		}

		order, err := domain.NewOrder(
			uint64(i+1),
			hookAddress,
			uint256.MustFromBig(new(big.Int).SetBytes(common.FromHex(orderRawData[1]))),
			uint256.MustFromBig(new(big.Int).SetBytes(common.FromHex(orderRawData[2]))),
			nil,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
