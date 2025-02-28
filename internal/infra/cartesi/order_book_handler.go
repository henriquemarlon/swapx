package cartesi

import (
	"encoding/hex"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/henriquemarlon/swapx/internal/infra/service"
	"github.com/henriquemarlon/swapx/internal/usecase"
	"github.com/henriquemarlon/swapx/pkg/coprocessor"
)

type OrderBookHandler struct {
	OrderRepository             domain.OrderRepository
	HookStorageServiceInterface service.HookStorageServiceInterface
}

func NewOrderHandler(orderRepository domain.OrderRepository, hookStorageServiceInterface service.HookStorageServiceInterface) *OrderBookHandler {
	return &OrderBookHandler{
		OrderRepository:             orderRepository,
		HookStorageServiceInterface: hookStorageServiceInterface,
	}
}

func (oh *OrderBookHandler) OrderBookHandler(input *coprocessor.AdvanceResponse) error {
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	inputArgs := abi.Arguments{
		{Type: uint256Type},
		{Type: uint256Type},
		{Type: uint256Type},
		{Type: uint256Type},
	}

	decodedData, err := hex.DecodeString(strings.TrimPrefix(input.Payload, "0x"))
	if err != nil {
		return err
	}

	values, err := inputArgs.Unpack(decodedData)
	if err != nil {
		return err
	}

	matchOrder := usecase.NewMatchOrderUseCase(
		oh.OrderRepository,
		oh.HookStorageServiceInterface,
	)
	res, err := matchOrder.Execute(&usecase.MatchOrderInputDTO{
		UnpackedArgs: values,
	}, input.Metadata)
	if err != nil {
		if err == domain.ErrNoMatch {
			log.Printf("No match found for order: %v", input.Metadata)
			return nil
		}
		return err
	}

	outputArgs := abi.Arguments{
		{Type: addressType},
		{Type: uint256Type},
		{Type: uint256Type},
	}

	encodedData, err := outputArgs.Pack(input.Metadata.MsgSender, res.BuyOrderId, res.SellOrderId)
	if err != nil {
		return err
	}

	coprocessor.SendNotice(&coprocessor.NoticeRequest{
		Payload: "0x" + common.Bytes2Hex(encodedData),
	})
	return nil
}
