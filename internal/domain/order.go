package domain

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	ErrNoMatch      = errors.New("no match found")
	ErrInvalidOrder = errors.New("invalid order")
	ErrOderNotFound = errors.New("order not found")
)

type OrderType string

const (
	OrderTypeBuy  OrderType = "buy"
	OrderTypeSell OrderType = "sell"
)

type OrderRepository interface {
	CreateOrder(order *Order) (*Order, error)
	FindAllOrders() ([]*Order, error)
	FindOrdersByType(orderType string) ([]*Order, error)
	FindOrderByTypeAndId(orderType string, id uint64) (*Order, error)
}

type Order struct {
	Id        uint64         `json:"id"`
	Account   common.Address `json:"account"`
	SqrtPrice *uint256.Int   `json:"sqrt_price"`
	Amount    *uint256.Int   `json:"amount"`
	Type      OrderType      `json:"type"`
}

func NewOrder(id uint64, account common.Address, sqrtPrice, amount *uint256.Int, orderType OrderType) (*Order, error) {
	order := &Order{
		Id:        id,
		Account:   account,
		SqrtPrice: sqrtPrice,
		Amount:    amount,
		Type:      orderType,
	}
	if err := order.Validate(); err != nil {
		return nil, err
	}
	return order, nil
}

func (o *Order) Validate() error {
	if o.Id == 0 || o.Account == (common.Address{}) || o.SqrtPrice.Sign() == 0 || o.Amount.Sign() == 0 {
		return ErrInvalidOrder
	}
	return nil
}
