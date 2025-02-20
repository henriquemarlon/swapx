package usecase

import (
	"errors"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/henriquemarlon/swapx/internal/infra/service"
	"github.com/henriquemarlon/swapx/pkg/coprocessor"
	"github.com/holiman/uint256"
)

var (
	infolog = log.New(os.Stderr, "[ info ]  ", log.Lshortfile)
)

type MatchOrderUseCase struct {
	OrderRepository domain.OrderRepository
}

type MatchOrderInputDTO struct {
	UnpackedArgs []interface{} `json:"unpacked_args"`
}

type MatchOrderOutputDTO struct {
	BuyOrderId  *big.Int `json:"buy_order_id"`
	SellOrderId *big.Int `json:"sell_order_id"`
}

func NewMatchOrderUseCase(orderRepository domain.OrderRepository) *MatchOrderUseCase {
	return &MatchOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (h *MatchOrderUseCase) Execute(input *MatchOrderInputDTO, metadata coprocessor.Metadata) (*MatchOrderOutputDTO, error) {
	amount, address, price, quantity, orderTypeBinary, err := validateAndExtractArgs(input.UnpackedArgs)
	if err != nil {
		return nil, err
	}

	var orderType domain.OrderType
	if orderTypeBinary.Cmp(big.NewInt(0)) == 0 {
		orderType = domain.OrderTypeBuy
	} else {
		orderType = domain.OrderTypeSell
	}

	order, err := domain.NewOrder(
		amount.Uint64(),
		address,
		uint256.MustFromBig(price),
		uint256.MustFromBig(quantity),
		orderType,
	)
	if err != nil {
		return nil, err
	}

	if _, err = h.OrderRepository.CreateOrder(order); err != nil {
		return nil, err
	}

	_, err = h.OrderRepository.FindAllOrders()
	if err != nil {
		return nil, err
	}

	factory := service.NewGioHandlerFactory("http://localhost:8080")
	handler, err := factory.NewGioHandler(0x27)
	if err != nil {
		return nil, err
	}

	_, err = handler.HandleStorageAt(
		common.HexToHash(metadata.BlockHash),
		metadata.MsgSender,
		common.BytesToHash(big.NewInt(8).Bytes()), // Change to a constant
	)
	if err != nil {
		return nil, err
	}

	return &MatchOrderOutputDTO{
		BuyOrderId:  big.NewInt(12),
		SellOrderId: big.NewInt(12),
	}, nil
}

func validateAndExtractArgs(args []interface{}) (*big.Int, common.Address, *big.Int, *big.Int, *big.Int, error) {
	if len(args) < 5 {
		return nil, common.Address{}, nil, nil, nil, errors.New("invalid input: UnpackedArgs must have at least 5 elements")
	}

	amount, ok := args[0].(*big.Int)
	if !ok {
		return nil, common.Address{}, nil, nil, nil, errors.New("invalid type for UnpackedArgs[0]: expected *big.Int")
	}

	address, ok := args[1].(common.Address)
	if !ok {
		return nil, common.Address{}, nil, nil, nil, errors.New("invalid type for UnpackedArgs[1]: expected common.Address")
	}

	price, ok := args[2].(*big.Int)
	if !ok {
		return nil, common.Address{}, nil, nil, nil, errors.New("invalid type for UnpackedArgs[2]: expected *big.Int")
	}

	quantity, ok := args[3].(*big.Int)
	if !ok {
		return nil, common.Address{}, nil, nil, nil, errors.New("invalid type for UnpackedArgs[3]: expected *big.Int")
	}

	orderTypeBinary, ok := args[4].(*big.Int)
	if !ok {
		return nil, common.Address{}, nil, nil, nil, errors.New("invalid type for UnpackedArgs[4]: expected *big.Int")
	}

	return amount, address, price, quantity, orderTypeBinary, nil
}
