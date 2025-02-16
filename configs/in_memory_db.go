package configs

import (
	"sync"

	"github.com/henriquemarlon/swapx/internal/domain"
)

type InMemoryDB struct {
	Orders map[uint64]*domain.Order
	Lock   *sync.RWMutex
}

func SetupInMemoryDB() (*InMemoryDB, error) {
	return &InMemoryDB{
		Orders: make(map[uint64]*domain.Order),
		Lock:   &sync.RWMutex{},
	}, nil
}
