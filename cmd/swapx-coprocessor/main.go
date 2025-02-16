package main

import (
	"encoding/json"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/swapx/configs"
	"github.com/henriquemarlon/swapx/internal/domain"
	"github.com/henriquemarlon/swapx/internal/infra/repository"
	"github.com/henriquemarlon/swapx/pkg/coprocessor"
)

var (
	infolog = log.New(os.Stderr, "[ info ]  ", log.Lshortfile)
	errlog  = log.New(os.Stderr, "[ error ] ", log.Lshortfile)
)

func Handler(response *coprocessor.AdvanceResponse) error {
	var newOrder *domain.Order
	var hookAddress common.Address

	// decode payload
	infolog.Println("Processing payload:", response)

	// call get storate GIO buyOrders (orderId, order)

	// call get storage GIO sellOrders (orderId, order)

	// setup database
	db, err := configs.SetupInMemoryDB()
	if err != nil {
		errlog.Panicln("Failed to setup database", "error", err)
	}

	// get all orders
	orderRepository := repository.NewOrderRepositoryInMemory(db)
	orders, err := orderRepository.FindAllOrders()
	if err != nil {
		errlog.Panicln("Failed to get all orders", "error", err)
	}
	infolog.Println("Database setup successful")

	// orderbook matching
	buyOrderId, sellOrderId, err := func(incomingOrder *domain.Order, orders []*domain.Order) (*big.Int, *big.Int, error) {
		return big.NewInt(12), big.NewInt(12), nil
	}(newOrder, orders)
	if err != nil {
		if err == domain.ErrNoMatch {
			infolog.Println("No match found")
			return nil
		}
		errlog.Panicln("Failed to match orders", "error", err)
	}

	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)

	arguments := abi.Arguments{
		{Type: addressType}, // "address"
		{Type: uint256Type}, // "uint256"
		{Type: uint256Type}, // "uint256"
	}

	encodedData, err := arguments.Pack(hookAddress, buyOrderId, sellOrderId)
	if err != nil {
		errlog.Panicln("Failed to encode ABI", "error", err)
	}
	
	res, err := coprocessor.SendNotice(&coprocessor.NoticeRequest{
		Payload: "0x" + common.Bytes2Hex(encodedData),
	})
	if err != nil {
		errlog.Panicln("Failed to send notice", "error", err)
	}
	infolog.Println("Notice sent", "status", res)
	return nil
}

func main() {
	finish := coprocessor.FinishRequest{Status: "accept"}
	for {
		infolog.Println("Sending finish")
		res, err := coprocessor.SendFinish(&finish)
		if err != nil {
			errlog.Panicln("Error: error making http request: ", err)
		}
		infolog.Println("Received finish status ", strconv.Itoa(res.StatusCode))

		if res.StatusCode == 202 {
			infolog.Println("No pending rollup request, trying again")
		} else {

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				errlog.Panicln("Error: could not read response body: ", err)
			}

			var finishResponse coprocessor.FinishResponse
			err = json.Unmarshal(resBody, &finishResponse)
			if err != nil {
				errlog.Panicln("Error: unmarshaling body:", err)
			}

			var rawPayload struct {
				Data string `json:"payload"`
			}
			if err := json.Unmarshal(finishResponse.Data, &rawPayload); err != nil {
				errlog.Println("Error unmarshaling payload:", err)
				finish.Status = "reject"
			}

			advanceResponse, err := coprocessor.EvmAdvanceParser(rawPayload.Data)
			if err != nil {
				errlog.Println(err)
				finish.Status = "reject"
			}

			err = Handler(&advanceResponse)
			if err != nil {
				errlog.Println(err)
				finish.Status = "reject"
			}
		}
	}
}
