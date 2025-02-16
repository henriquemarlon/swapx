package usecase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/holiman/uint256"
)

type CreateOrderInputDTO struct {
	Id        uint64         `json:"id"`
	Account   common.Address `json:"account"`
	SqrtPrice *uint256.Int   `json:"sqrt_price"`
	Amount    *uint256.Int   `json:"amount"`
	Type      string         `json:"type"`
}

type CreateOrderUseCase struct {
	OrderRepository domain.OrderRepository
}

func NewCreateOrderUseCase(orderRepository domain.OrderRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (u *CreateOrderUseCase) Execute(input *CreateOrderInputDTO) (*FindOrderOutputDTO, error) {
	order, err := domain.NewOrder(
		input.Id,
		input.Account,
		input.SqrtPrice,
		input.Amount,
		domain.OrderType(input.Type),
	)
	if err != nil {
		return nil, err
	}

	res, err := u.OrderRepository.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	return &FindOrderOutputDTO{
		Id:        res.Id,
		Account:   res.Account,
		SqrtPrice: res.SqrtPrice,
		Amount:    res.Amount,
		Type:      string(res.Type),
	}, nil
}
