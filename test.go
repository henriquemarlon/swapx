package main

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// ABI da função EvmAdvance exatamente como no contrato Solidity
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
				{"name": "payload", "type": "bytes"}
			],
			"outputs": []
		}
	]`

	// Parseia o ABI da função
	parsedABI, err := abi.JSON(strings.NewReader(evmAdvanceABI))
	if err != nil {
		log.Fatalf("Erro ao parsear ABI: %v", err)
	}

	// Verificar se a função EvmAdvance está no ABI carregado
	method, exists := parsedABI.Methods["EvmAdvance"]
	if !exists {
		log.Fatalf("Método EvmAdvance não encontrado no ABI")
	}

	// Payload recebido em formato hexadecimal (exemplo: substitua pelo payload real)
	payload := common.FromHex("0x477273f60000000000000000000000000000000000000000000000000000000000000001000000000000000000000000358aa13c52544eccef6b0add0f801012adad5ee30000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4509a90371e07e93332ca15363c67b348e0d0ec85c27be20a908bfeab7a185ebc00000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000067b186a400000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000012")

	// Remover os primeiros 4 bytes do selector da função
	if len(payload) < 4 {
		log.Fatalf("Payload muito curto para conter um selector válido")
	}
	payload = payload[4:]

	// Decodifica os argumentos do payload
	args, err := method.Inputs.Unpack(payload)
	if err != nil {
		log.Fatalf("Erro ao decodificar payload: %v", err)
	}

	// Extrair os valores decodificados
	chainId := args[0].(*big.Int)
	taskManager := args[1].(common.Address)
	msgSender := args[2].(common.Address)
	blockHash := args[3].([32]byte)
	blockNumber := args[4].(*big.Int)
	blockTimestamp := args[5].(*big.Int)
	decodedPayload := args[6].([]byte)

	// Exibir os valores extraídos
	fmt.Printf("Chain ID: %s\n", chainId.String())
	fmt.Printf("Task Manager: %s\n", taskManager.Hex())
	fmt.Printf("Msg Sender: %s\n", msgSender.Hex())
	fmt.Printf("Block Hash: 0x%x\n", blockHash)
	fmt.Printf("Block Number: %s\n", blockNumber.String())
	fmt.Printf("Block Timestamp: %s\n", blockTimestamp.String())
	fmt.Printf("Payload (decodificado): %v\n", decodedPayload)
}
