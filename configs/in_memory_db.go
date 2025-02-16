package configs

import (
	"sync"

	"github.com/henriquemarlon/swapx/internal/domain"
)

type InMemoryDB struct {
	Orders map[uint]*domain.Order
	Lock   *sync.RWMutex
}

func SetupInMemoryDB() (*InMemoryDB, error) {
	return &InMemoryDB{
		Orders: make(map[uint]*domain.Order),
		Lock:   &sync.RWMutex{},
	}, nil
}
