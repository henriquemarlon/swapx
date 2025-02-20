package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type GioGetStorage struct {
	Domain  uint16 `json:"domain"`
	BaseUrl string `json:"base_url"`
}

func (h *GioGetStorage) HandleStorageAt(blockHash common.Hash, address common.Address, slot common.Hash) (*GioResponse, error) {
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

	req, err := http.NewRequest("POST", h.BaseUrl+"/gio", bytes.NewBuffer(encodedData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", res.StatusCode, string(body))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var gioResp GioResponse
	if err := json.Unmarshal(body, &gioResp); err != nil {
		return nil, err
	}

	return &gioResp, nil
}
