package configs

import (
	"sync"

	"github.com/henriquemarlon/swapx/internal/domain"
)

type InMemoryDB struct {
	BuyOrders  map[uint64]*domain.Order
	SellOrders map[uint64]*domain.Order
	Mutex       *sync.RWMutex
}

func SetupInMemoryDB() (*InMemoryDB, error) {
	return &InMemoryDB{
		BuyOrders:  make(map[uint64]*domain.Order),
		SellOrders: make(map[uint64]*domain.Order),
		Mutex:       &sync.RWMutex{},
	}, nil
}
