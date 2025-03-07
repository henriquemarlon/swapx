package usecase

import (
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/henriquemarlon/swapx/internal/infra/service"
	"github.com/henriquemarlon/swapx/pkg/coprocessor"
	"github.com/holiman/uint256"
)

const (
	BUY_ORDERS_STORAGE_SLOT  = 8
	SELL_ORDERS_STORAGE_SLOT = 9
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
		index.Uint64() + 1,
		metadata.MsgSender,
		uint256.MustFromBig(price),
		uint256.MustFromBig(quantity),
		&orderType,
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
	)
	if err != nil {
		if err != service.ErrNoOrdersFound {
			return nil, errors.New("cannot match order: no buy orders found")
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
	)
	if err != nil {
		if err != service.ErrNoOrdersFound {
			return nil, errors.New("cannot match order: no sell orders found")
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

	bids, err := h.OrderRepository.FindOrdersByType(string(domain.OrderTypeBuy))
	if err != nil {
		log.Printf("Error fetching all orders: %v", err)
		return nil, err
	}
	for _, bid := range bids {
		orderBook.Bids.Push(bid)
	}

	asks, err := h.OrderRepository.FindOrdersByType(string(domain.OrderTypeSell))
	if err != nil {
		log.Printf("Error fetching all orders: %v", err)
		return nil, err
	}
	for _, ask := range asks {
		orderBook.Asks.Push(ask)
	}

	trades, err := orderBook.MatchOrders()
	if err != nil {
		return nil, err
	}

	return &MatchOrdersOutputDTO{
		Trades: trades,
	}, nil
}
