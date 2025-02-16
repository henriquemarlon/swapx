package usecase

import (
	"github.com/henriquemarlon/swapx/internal/domain"
)

type FindOrderByTypeAndIdInputDTO struct {
	Type string `json:"type"`
	Id   uint64 `json:"id"`
}

type FindOrderByTypeAndIdUsecase struct {
	OrderRepository domain.OrderRepository
}

func NewFindOrderByTypeAndIdUseCase(orderRepository domain.OrderRepository) *FindOrderByTypeAndIdUsecase {
	return &FindOrderByTypeAndIdUsecase{
		OrderRepository: orderRepository,
	}
}

func (s *FindOrderByTypeAndIdUsecase) Execute(input *FindOrderByTypeAndIdInputDTO) (*FindOrderOutputDTO, error) {
	res, err := s.OrderRepository.FindOrderByTypeAndId(input.Type, input.Id)
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
