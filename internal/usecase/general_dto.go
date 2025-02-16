package usecase

import "github.com/holiman/uint256"

type FindOrderOutputDTO struct {
	Id        uint         `json:"id"`
	Account   string       `json:"account"`
	SqrtPrice *uint256.Int `json:"sqrt_price"`
	Amount    *uint256.Int `json:"amount"`
	Type      string       `json:"type"`
}
