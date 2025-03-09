package repository

import (
	"sync"

	"github.com/henriquemarlon/swapx/configs"
	"github.com/henriquemarlon/swapx/internal/domain"
)

type OrderRepositoryInMemory struct {
	BuyOrders  map[uint64]*domain.Order
	SellOrders map[uint64]*domain.Order
	Mutex      *sync.RWMutex
}

func NewOrderRepositoryInMemory(db *configs.InMemoryDB) *OrderRepositoryInMemory {
	return &OrderRepositoryInMemory{
		BuyOrders:  db.BuyOrders,
		SellOrders: db.SellOrders,
		Mutex:      db.Mutex,
	}
}

func (r *OrderRepositoryInMemory) CreateOrder(order *domain.Order) (*domain.Order, error) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	orderMap := r.getOrderMap((*string)(order.Type))
	if _, exists := orderMap[order.Id]; exists {
		return nil, domain.ErrOrderAlreadyExists
	}

	orderMap[order.Id] = order
	return order, nil
}

func (r *OrderRepositoryInMemory) FindAllOrders() ([]*domain.Order, error) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	var orders []*domain.Order
	for _, order := range r.BuyOrders {
		orders = append(orders, order)
	}
	for _, order := range r.SellOrders {
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return nil, domain.ErrNoOrdersFound
	}
	return orders, nil
}

func (r *OrderRepositoryInMemory) FindOrderById(orderType string, id uint64) (*domain.Order, error) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	orderMap := r.getOrderMap(&orderType)
	order, exists := orderMap[id]
	if !exists {
		return nil, domain.ErrOrderNotFound
	}
	return order, nil
}

func (r *OrderRepositoryInMemory) FindOrdersByType(orderType string) ([]*domain.Order, error) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	orderMap := r.getOrderMap(&orderType)
	var orders []*domain.Order
	for _, order := range orderMap {
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return nil, domain.ErrNoOrdersFound
	}
	return orders, nil
}

func (r *OrderRepositoryInMemory) getOrderMap(orderType *string) map[uint64]*domain.Order {
	if *orderType == string(domain.OrderTypeBuy) {
		return r.BuyOrders
	}
	return r.SellOrders
}
