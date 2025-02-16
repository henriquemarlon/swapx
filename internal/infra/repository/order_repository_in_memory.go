package repository

import (
	"log"
	"os"
	"sync"

	"github.com/henriquemarlon/swapx/configs"
	"github.com/henriquemarlon/swapx/internal/domain"
)

var (
	infolog = log.New(os.Stderr, "[ info ]  ", log.Lshortfile)
)

type OrderRepositoryInMemory struct {
	db    map[uint64]*domain.Order
	mutex *sync.RWMutex
}

func NewOrderRepositoryInMemory(db *configs.InMemoryDB) *OrderRepositoryInMemory {
	return &OrderRepositoryInMemory{
		db:    db.Orders,
		mutex: db.Lock,
	}
}

func (r *OrderRepositoryInMemory) CreateOrder(order *domain.Order) (*domain.Order, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.db[order.Id] = order
	infolog.Println("Order created:", order)
	return order, nil
}

func (r *OrderRepositoryInMemory) FindAllOrders() ([]*domain.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var orders []*domain.Order
	for _, order := range r.db {
		orders = append(orders, order)
		infolog.Println("Found", "order", order)
	}
	return orders, nil
}

func (r *OrderRepositoryInMemory) FindOrderByTypeAndId(orderType string, id uint64) (*domain.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	order, exists := r.db[id]
	if !exists || string(order.Type) != orderType {
		return nil, domain.ErrOderNotFound
	}
	infolog.Println("Found", "order", order)
	return order, nil
}

func (r *OrderRepositoryInMemory) FindOrdersByType(orderType string) ([]*domain.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var filteredOrders []*domain.Order
	for _, order := range r.db {
		if string(order.Type) == orderType {
			filteredOrders = append(filteredOrders, order)
			infolog.Println("Found", "order", order)
		}
	}
	return filteredOrders, nil
}
