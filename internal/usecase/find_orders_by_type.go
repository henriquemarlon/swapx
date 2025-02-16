package usecase

import (
	"github.com/henriquemarlon/swapx/internal/domain"
)

type FindOrdersByTypeInputDTO struct {
	Type string `json:"type"`
}

type FindOrdersByTypeOutputDTO []*FindOrderOutputDTO

type FindOrdersByTypeUsecase struct {
	OrderRepository domain.OrderRepository
}

func NewFindOrdersByTypeUseCase(orderRepository domain.OrderRepository) *FindOrdersByTypeUsecase {
	return &FindOrdersByTypeUsecase{
		OrderRepository: orderRepository,
	}
}

func (s *FindOrdersByTypeUsecase) Execute(input *FindOrdersByTypeInputDTO) (FindOrdersByTypeOutputDTO, error) {
	res, err := s.OrderRepository.FindOrdersByType(input.Type)
	if err != nil {
		return nil, err
	}
	var output FindOrdersByTypeOutputDTO
	for _, order := range res {
		dto := &FindOrderOutputDTO{
			Id:        order.Id,
			Account:   order.Account,
			SqrtPrice: order.SqrtPrice,
			Amount:    order.Amount,
			Type:      string(order.Type),
		}
		output = append(output, dto)
	}
	return output, nil
}
