//go:build wireinject
// +build wireinject

package root

import (
	"github.com/google/wire"
	"github.com/henriquemarlon/swapx/configs"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/henriquemarlon/swapx/internal/infra/cartesi"
	"github.com/henriquemarlon/swapx/internal/infra/repository"
)

var setOrderRepositoryDependency = wire.NewSet(
	repository.NewOrderRepositoryInMemory,
	wire.Bind(new(domain.OrderRepository), new(*repository.OrderRepositoryInMemory)),
)

var setOrderHandler = wire.NewSet(
	cartesi.NewOrderHandler,
)

func NewOrderHandler(db *configs.InMemoryDB) (*cartesi.OrderBookHandler, error) {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderHandler,
	)
	return nil, nil
}
