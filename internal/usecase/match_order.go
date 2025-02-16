package usecase

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/henriquemarlon/swapx/pkg/coprocessor"
	"github.com/holiman/uint256"
)

type MatchOrderUseCase struct {
	OrderRepository domain.OrderRepository
}

type MatchOrderInputDTO struct {
	UnpackedArgs []interface{} `json:"unpacked_args"`
}

type MatchOrderOutputDTO struct {
	BuyOrderId  *big.Int `json:"buy_order_id"`
	SellOrderId *big.Int `json:"sell_order_id"`
}

func NewMatchOrderUseCase(orderRepository domain.OrderRepository) *MatchOrderUseCase {
	return &MatchOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (h *MatchOrderUseCase) Execute(input *MatchOrderInputDTO, metadata coprocessor.Metadata) (*MatchOrderOutputDTO, error) {
	var orderType domain.OrderType
	if input.UnpackedArgs[4].(*big.Int).Cmp(big.NewInt(0)) == 0 {
		orderType = domain.OrderTypeBuy
	} else {
		orderType = domain.OrderTypeSell
	}

	order, err := domain.NewOrder(
		input.UnpackedArgs[0].(*big.Int).Uint64(),
		input.UnpackedArgs[1].(common.Address),
		uint256.MustFromBig(input.UnpackedArgs[2].(*big.Int)),
		uint256.MustFromBig(input.UnpackedArgs[3].(*big.Int)),
		orderType,
	)
	if err != nil {
		return nil, err
	}

	if _, err = h.OrderRepository.CreateOrder(order); err != nil {
		return nil, err
	}

	_, err = h.OrderRepository.FindAllOrders()
	if err != nil {
		return nil, err
	}

	// call get storate GIO buyOrders (orderId, order)([]bytes)

	// call get storage GIO sellOrders (orderId, order)

	// orderbook match making

	return &MatchOrderOutputDTO{
		BuyOrderId:  big.NewInt(12),
		SellOrderId: big.NewInt(12),
	}, nil
}
