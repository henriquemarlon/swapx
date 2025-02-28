package gio

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
	Handle(blockHash common.Hash, address common.Address, slot common.Hash) (*GioResponse, error)
}
