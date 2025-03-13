package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/henriquemarlon/swapx/internal/infra/service"
	"github.com/henriquemarlon/swapx/pkg/coprocessor"
	"github.com/holiman/uint256"
)

const (
	BUY_ORDERS_STORAGE_SLOT         = 8
	BUY_ORDERS_STATUS_STORAGE_SLOT  = 6
	SELL_ORDERS_STORAGE_SLOT        = 9
	SELL_ORDERS_STATUS_STORAGE_SLOT = 7
)

type MatchOrdersUseCase struct {
	OrderRepository     domain.OrderRepository
	HookContractService service.OrderStorageServiceInterface
}

type MatchOrdersInputDTO struct {
	UnpackedArgs []interface{} `json:"unpacked_args"`
}

type MatchOrdersOutputDTO struct {
	Trades []*domain.Trade `json:"trades"`
}

func NewMatchOrdersUseCase(orderRepository domain.OrderRepository, hookContractService service.OrderStorageServiceInterface) *MatchOrdersUseCase {
	return &MatchOrdersUseCase{
		OrderRepository:     orderRepository,
		HookContractService: hookContractService,
	}
}

func (h *MatchOrdersUseCase) Execute(input *MatchOrdersInputDTO, metadata coprocessor.Metadata) (*MatchOrdersOutputDTO, error) {

	// -----------------------------------------------------------------------------
	// Validate input
	// -----------------------------------------------------------------------------

	if len(input.UnpackedArgs) < 4 {
		return nil, errors.New("invalid input: UnpackedArgs must have at least 4 elements")
	}

	index, ok := input.UnpackedArgs[0].(*big.Int)
	if !ok {
		return nil, errors.New("invalid type for UnpackedArgs[0]: expected *big.Int")
	}

	price, ok := input.UnpackedArgs[1].(*big.Int)
	if !ok {
		return nil, errors.New("invalid type for UnpackedArgs[1]: expected *big.Int")
	}

	quantity, ok := input.UnpackedArgs[2].(*big.Int)
	if !ok {
		return nil, errors.New("invalid type for UnpackedArgs[2]: expected *big.Int")
	}

	orderTypeBinary, ok := input.UnpackedArgs[3].(*big.Int)
	if !ok {
		return nil, errors.New("invalid type for UnpackedArgs[3]: expected *big.Int")
	}

	// -----------------------------------------------------------------------------
	// Create incoming order
	// -----------------------------------------------------------------------------

	var orderType domain.OrderType
	if orderTypeBinary.Cmp(big.NewInt(0)) == 0 {
		orderType = domain.OrderTypeBuy
	} else {
		orderType = domain.OrderTypeSell
	}

	order, err := domain.NewOrder(
		index.Uint64(), // The index here comes from the order array length at the time of the swap call, which is (orderIndex + 1).
		metadata.MsgSender,
		uint256.MustFromBig(price),
		uint256.MustFromBig(quantity),
		uint256.MustFromBig(big.NewInt(0)),
		&orderType,
		&domain.OrderNotCancelledOrFulfilled,
	)
	if err != nil {
		return nil, err
	}

	if _, err = h.OrderRepository.CreateOrder(order); err != nil {
		return nil, err
	}

	// -----------------------------------------------------------------------------
	// Find all previous orders ( Base layer access )
	// -----------------------------------------------------------------------------

	buyOrders, err := h.HookContractService.FindOrdersBySlot(
		metadata.MsgSender,
		common.HexToHash(metadata.BlockHash),
		common.BigToHash(big.NewInt(BUY_ORDERS_STORAGE_SLOT)),
		common.BigToHash(big.NewInt(BUY_ORDERS_STATUS_STORAGE_SLOT)),
	)
	if err != nil {
		if err == domain.ErrNoOrdersFound {
			return nil, fmt.Errorf("%v with type buy", err)
		}
		buyOrders = nil
	}

	for _, buyOrder := range buyOrders {
		buyOrder.Type = &domain.OrderTypeBuy
		if _, err := h.OrderRepository.CreateOrder(buyOrder); err != nil {
			return nil, err
		}
	}

	sellOrders, err := h.HookContractService.FindOrdersBySlot(
		metadata.MsgSender,
		common.HexToHash(metadata.BlockHash),
		common.BigToHash(big.NewInt(SELL_ORDERS_STORAGE_SLOT)),
		common.BigToHash(big.NewInt(SELL_ORDERS_STATUS_STORAGE_SLOT)),
	)
	if err != nil {
		if err == domain.ErrNoOrdersFound {
			return nil, fmt.Errorf("%v with type sell", err)
		}
		sellOrders = nil
	}

	for _, sellOrder := range sellOrders {
		sellOrder.Type = &domain.OrderTypeSell
		if _, err := h.OrderRepository.CreateOrder(sellOrder); err != nil {
			return nil, err
		}
	}

	// -----------------------------------------------------------------------------
	// Match orders
	// -----------------------------------------------------------------------------

	orderBook := domain.NewOrderBook()

	bids, err := h.OrderRepository.FindOrdersByTypeAndStatus(domain.OrderTypeBuy, domain.OrderNotCancelledOrFulfilled)
	if err != nil {
		if err == domain.ErrNoOrdersFound {
			return nil, fmt.Errorf("%v with type buy in memory", err)
		}
		return nil, err
	}
	for _, bid := range bids {
		orderBook.Bids.Push(bid)
	}

	asks, err := h.OrderRepository.FindOrdersByTypeAndStatus(domain.OrderTypeSell, domain.OrderNotCancelledOrFulfilled)
	if err != nil {
		if err == domain.ErrNoOrdersFound {
			return nil, fmt.Errorf("%v with type sell in memory", err)
		}
		return nil, err
	}
	for _, ask := range asks {
		orderBook.Asks.Push(ask)
	}

	orders, err := h.OrderRepository.FindAllOrders()
	if err != nil {
		return nil, err
	}

	ordersBytes, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}
	slog.Info("Current state before match", "info", string(ordersBytes))

	trades, err := orderBook.MatchOrders()
	if err != nil {
		return nil, err
	}

	tradesBytes, err := json.Marshal(trades)
	if err != nil {
		return nil, err
	}
	slog.Info("Selected trades", "info", string(tradesBytes))

	return &MatchOrdersOutputDTO{
		Trades: trades,
	}, nil
}
