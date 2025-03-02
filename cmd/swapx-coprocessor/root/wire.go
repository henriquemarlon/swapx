//go:build wireinject
// +build wireinject

package root

import (
	"github.com/google/wire"
	"github.com/henriquemarlon/swapx/configs"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/henriquemarlon/swapx/internal/infra/cartesi"
	"github.com/henriquemarlon/swapx/internal/infra/repository"
	"github.com/henriquemarlon/swapx/internal/infra/service"
	"github.com/henriquemarlon/swapx/pkg/gio"
)

var setHookStorageService = wire.NewSet(
	service.NewOrderStorageService,
	wire.Bind(new(service.OrderStorageServiceInterface), new(*service.OrderStorageService)),
)

var setGioHandlerFactory = wire.NewSet(
	gio.NewGioHandlerFactory,
)

var setOrderRepositoryDependency = wire.NewSet(
	repository.NewOrderRepositoryInMemory,
	wire.Bind(new(domain.OrderRepository), new(*repository.OrderRepositoryInMemory)),
)

var setMatchOrdersHandler = wire.NewSet(
	cartesi.NewMatchOrdersHandler,
)

func NewMatchOrdersHandler(db *configs.InMemoryDB, rollupServerUrl string) (*cartesi.MatchOrdersHandler, error) {
	wire.Build(
		setOrderRepositoryDependency,
		setGioHandlerFactory,
		setHookStorageService,
		setMatchOrdersHandler,
	)
	return &cartesi.MatchOrdersHandler{}, nil
}