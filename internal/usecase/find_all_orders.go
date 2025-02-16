package usecase

import (
	"github.com/henriquemarlon/swapx/internal/domain"
)

type FindAllOrdersOutputDTO []*FindOrderOutputDTO

type FindAllOrdersUsecase struct {
	OrderRepository domain.OrderRepository
}

func NewFindAllOrdersUseCase(orderRepository domain.OrderRepository) *FindAllOrdersUsecase {
	return &FindAllOrdersUsecase{
		OrderRepository: orderRepository,
	}
}

func (s *FindAllOrdersUsecase) Execute() (FindAllOrdersOutputDTO, error) {
	res, err := s.OrderRepository.FindAllOrders()
	if err != nil {
		return nil, err
	}
	var output FindAllOrdersOutputDTO
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
