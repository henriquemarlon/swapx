package service

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type GioGetStorage struct {
	Domain  uint16 `json:"domain"`
	BaseUrl string `json:"base_url"`
}

func (h *GioGetStorage) HandleStorageAt(blockHash common.Hash, address common.Address, slot common.Hash) (*GioResponse, error) {
	log.Printf("Handling storage at block %s, address %s, slot %s\n", blockHash.Hex(), address.Hex(), slot.Hex())
	client := &http.Client{}

	addressType, _ := abi.NewType("address", "", nil)
	bytes32Type, _ := abi.NewType("bytes32", "", nil)

	outputArgs := abi.Arguments{
		{Type: bytes32Type},
		{Type: addressType},
		{Type: bytes32Type},
	}

	if len(blockHash) != 32 {
		return nil, errors.New("invalid block hash format")
	}

	encodedData, err := outputArgs.Pack(blockHash, address, slot)
	if err != nil {
		return nil, err
	}

	log.Printf("Encoded data: %v\n", encodedData)
	hexEncoded := hex.EncodeToString(encodedData)
	reqBody, err := json.Marshal(GioRequest{
		Domain: 0x27,
		Id:     "0x" + hexEncoded,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:5004/gio", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Response status: %d, body: %s\n", res.StatusCode, string(body))
		return nil, errors.New("unexpected status code: " + res.Status + ", response: " + string(body))
	}

	var gioResponse *GioResponse
	if err := json.Unmarshal(body, &gioResponse); err != nil {
		return nil, errors.New("invalid JSON response format: " + err.Error())
	}

	return gioResponse, nil
}
