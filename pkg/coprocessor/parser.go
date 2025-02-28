package coprocessor

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const evmAdvanceABI = `[
	{
		"name": "EvmAdvance",
		"type": "function",
		"stateMutability": "nonpayable",
		"inputs": [
			{"name": "chainId", "type": "uint256"},
			{"name": "taskManager", "type": "address"},
			{"name": "msgSender", "type": "address"},
			{"name": "blockHash", "type": "bytes32"},
			{"name": "blockNumber", "type": "uint256"},
			{"name": "blockTimestamp", "type": "uint256"},
			{"name": "prevRandao", "type": "uint256"},
			{"name": "payload", "type": "bytes"}
		],
		"outputs": []
	}
]`

func EvmAdvanceParser(hexInput string) (AdvanceResponse, error) {
	var response AdvanceResponse

	parsedABI, err := abi.JSON(strings.NewReader(evmAdvanceABI))
	if err != nil {
		return response, fmt.Errorf("error parsing ABI: %v", err)
	}

	method, exists := parsedABI.Methods["EvmAdvance"]
	if !exists {
		return response, fmt.Errorf("method EvmAdvance not found in ABI")
	}

	payload := common.FromHex(hexInput)
	if len(payload) < 4 {
		return response, fmt.Errorf("payload too short to contain a valid selector")
	}
	payload = payload[4:]

	args, err := method.Inputs.Unpack(payload)
	if err != nil {
		return response, fmt.Errorf("error decoding payload: %v", err)
	}

	chainId := args[0].(*big.Int).Uint64()
	taskManager := args[1].(common.Address)
	msgSender := args[2].(common.Address)
	blockHash := fmt.Sprintf("0x%x", args[3].([32]byte))
	blockNumber := args[4].(*big.Int).Uint64()
	blockTimestamp := args[5].(*big.Int).Uint64()
	prevRandao := args[6].(*big.Int).String()
	decodedPayload := fmt.Sprintf("0x%x", args[7].([]byte))

	response.Metadata.ChainId = chainId
	response.Metadata.TaskManager = taskManager
	response.Metadata.MsgSender = msgSender
	response.Metadata.BlockHash = blockHash
	response.Metadata.BlockNumber = blockNumber
	response.Metadata.Timestamp = blockTimestamp
	response.Metadata.PrevRandao = prevRandao
	response.Payload = decodedPayload

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return response, fmt.Errorf("error marshaling response: %v", err)
	}
	log.Printf("Advance response: %s", string(jsonBytes))

	return response, nil
}
