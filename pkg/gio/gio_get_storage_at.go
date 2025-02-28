package gio

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type GioGetStorage struct {
	Domain  uint16 `json:"domain"`
	BaseUrl string `json:"base_url"`
}

func NewGioGetStorage(baseUrl string, domain uint16) *GioGetStorage {
	return &GioGetStorage{
		Domain:  domain,
		BaseUrl: baseUrl,
	}
}

func (h *GioGetStorage) Handle(blockHash common.Hash, address common.Address, slot common.Hash) (*GioResponse, error) {
	log.Printf("Handling storage at block %s, address %s, slot %s\n", blockHash.Hex(), address.Hex(), slot.Hex())

	hexEncoded := append(blockHash[:], address[:]...)
	hexEncoded = append(hexEncoded, slot[:]...)

	encodedData := hex.EncodeToString(hexEncoded)

	reqBody, err := json.Marshal(GioRequest{
		Domain: 0x27,
		Id:     "0x" + encodedData,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Request body: %s\n", string(reqBody))

	req, err := http.NewRequest("POST", h.BaseUrl + "/gio", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusAccepted {
		log.Printf("Response status: %d, body: %s\n", res.StatusCode, string(body))
		return nil, errors.New("unexpected status code: " + res.Status + ", response: " + string(body))
	}

	var gioResponse GioResponse
	if err := json.Unmarshal(body, &gioResponse); err != nil {
		return nil, errors.New("invalid JSON response format: " + err.Error())
	}

	return &gioResponse, nil
}
