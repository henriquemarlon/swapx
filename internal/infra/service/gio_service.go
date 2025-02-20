package service

import "github.com/ethereum/go-ethereum/common"

type GioRequest struct {
	Domain uint16 `json:"domain"`
	Id     string `json:"id"`
}

type GioResponse struct {
	ResponseCode uint16 `json:"response_code"`
	Response     string `json:"response"`
}

type GioHandler interface {
	HandleStorageAt(slot common.Hash, address common.Address, blockHash common.Hash) (*GioResponse, error)
}
