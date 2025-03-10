package domain

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	ErrInvalidOrder       = errors.New("invalid order")
	ErrOrderNotFound      = errors.New("order not found")
	ErrNoOrdersFound      = errors.New("no orders found")
	ErrOrderAlreadyExists = errors.New("order with this id already exists")
)

type OrderType string

var (
	OrderTypeBuy  OrderType = "buy"
	OrderTypeSell OrderType = "sell"
)

type OrderStatus string

var (
	OrderStatusOpen              OrderStatus = "open"
	OrderStatusFulFilledOrClosed OrderStatus = "fulfilled_or_closed"
)

type OrderRepository interface {
	FindAllOrders() ([]*Order, error)
	CreateOrder(order *Order) (*Order, error)
	FindOrdersByType(orderType OrderType) ([]*Order, error)
	FindOrderById(orderType OrderType, id uint64) (*Order, error)
	FindOrdersByTypeAndStatus(orderType OrderType, orderStatus OrderStatus) ([]*Order, error)
}

type Order struct {
	Id        uint64         `json:"id"`
	Hook      common.Address `json:"hook"`
	SqrtPrice *uint256.Int   `json:"sqrt_price"`
	Amount    *uint256.Int   `json:"amount"`
	Type      *OrderType     `json:"type"`
	Status    *OrderStatus   `json:"status"`
}

func NewOrder(id uint64, hook common.Address, sqrtPrice, amount *uint256.Int, orderType *OrderType, orderStatus *OrderStatus) (*Order, error) {
	order := &Order{
		Id:        id,
		Hook:      hook,
		SqrtPrice: sqrtPrice,
		Amount:    amount,
		Type:      orderType,
		Status:    orderStatus,
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
