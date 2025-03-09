package domain

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	ErrInvalidOrder = errors.New("invalid order")
	ErrOrderNotFound = errors.New("order not found")
	ErrNoOrdersFound     = errors.New("no orders found")
	ErrOrderAlreadyExists = errors.New("order with this id already exists")
)

type OrderType string

var (
	OrderTypeBuy  OrderType = "buy"
	OrderTypeSell OrderType = "sell"
)

type OrderRepository interface {
	CreateOrder(order *Order) (*Order, error)
	FindAllOrders() ([]*Order, error)
	FindOrdersByType(orderType string) ([]*Order, error)
	FindOrderById(orderType string, id uint64) (*Order, error)
}

type Order struct {
	Id        uint64         `json:"id"`
	Hook      common.Address `json:"hook"`
	SqrtPrice *uint256.Int   `json:"sqrt_price"`
	Amount    *uint256.Int   `json:"amount"`
	Type      *OrderType     `json:"type"`
}

func NewOrder(id uint64, hook common.Address, sqrtPrice, amount *uint256.Int, orderType *OrderType) (*Order, error) {
	order := &Order{
		Id:        id,
		Hook:      hook,
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
	if o.Id == 0 || o.Hook == (common.Address{}) || o.SqrtPrice.Sign() == 0 || o.Amount.Sign() == 0 {
		return ErrInvalidOrder
	}
	return nil
}
