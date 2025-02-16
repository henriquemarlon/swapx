package usecase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type FindOrderOutputDTO struct {
	Id        uint64         `json:"id"`
	Account   common.Address `json:"account"`
	SqrtPrice *uint256.Int   `json:"sqrt_price"`
	Amount    *uint256.Int   `json:"amount"`
	Type      string         `json:"type"`
}
