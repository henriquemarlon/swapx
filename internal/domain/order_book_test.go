package domain

import (
	"container/heap"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func setupOrderBook(bids, asks []*Order) *OrderBook {
	orderBook := NewOrderBook()
	for _, bid := range bids {
		heap.Push(orderBook.Bids, bid)
	}
	for _, ask := range asks {
		heap.Push(orderBook.Asks, ask)
	}
	return orderBook
}

var testHook = common.HexToAddress("0x1")

func TestBidFullyMatchedBySingleAsk(t *testing.T) {
	bids := []*Order{
		{
			Id:            1,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(100),
			Amount:        uint256.NewInt(50),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeBuy,
		},
	}
	asks := []*Order{
		{
			Id:            2,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(90),
			Amount:        uint256.NewInt(50),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeSell,
		},
	}
	orderBook := setupOrderBook(bids, asks)
	expectedTrades := []*Trade{{BidId: 1, AskId: 2}}
	trades, err := orderBook.MatchOrders()

	assert.NoError(t, err)
	assert.NotNil(t, trades)
	assert.Equal(t, expectedTrades, trades)
}

func TestBidFullyMatchedByMultipleAsks(t *testing.T) {
	bids := []*Order{
		{
			Id:            1,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(100),
			Amount:        uint256.NewInt(100),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeBuy,
		},
	}
	asks := []*Order{
		{
			Id:            2,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(90),
			Amount:        uint256.NewInt(40),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeSell,
		},
		{
			Id:            3,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(85),
			Amount:        uint256.NewInt(60),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeSell,
		},
	}
	orderBook := setupOrderBook(bids, asks)
	expectedTrades := []*Trade{{BidId: 1, AskId: 3}, {BidId: 1, AskId: 2}}
	trades, err := orderBook.MatchOrders()

	assert.NoError(t, err)
	assert.NotNil(t, trades)
	assert.Equal(t, expectedTrades, trades)
}

func TestBidPartiallyMatched(t *testing.T) {
	bids := []*Order{
		{
			Id:            1,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(100),
			Amount:        uint256.NewInt(80),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeBuy,
		},
	}
	asks := []*Order{
		{
			Id:            2,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(90),
			Amount:        uint256.NewInt(50),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeSell,
		},
		{
			Id:            3,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(100),
			Amount:        uint256.NewInt(40),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeSell,
		},
	}
	orderBook := setupOrderBook(bids, asks)
	expectedTrades := []*Trade{{BidId: 1, AskId: 2}, {BidId: 1, AskId: 3}}
	trades, err := orderBook.MatchOrders()

	assert.NoError(t, err)
	assert.NotNil(t, trades)
	assert.Equal(t, expectedTrades, trades)
}

func TestAskFullyMatchedBySingleBid(t *testing.T) {
	bids := []*Order{
		{
			Id:            1,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(100),
			Amount:        uint256.NewInt(50),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeBuy,
		},
	}
	asks := []*Order{
		{
			Id:            2,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(100),
			Amount:        uint256.NewInt(50),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeSell,
		},
	}
	orderBook := setupOrderBook(bids, asks)
	expectedTrades := []*Trade{{BidId: 1, AskId: 2}}
	trades, err := orderBook.MatchOrders()

	assert.NoError(t, err)
	assert.NotNil(t, trades)
	assert.Equal(t, expectedTrades, trades)
}

func TestAskFullyMatchedByMultipleBids(t *testing.T) {
	bids := []*Order{
		{
			Id:            1,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(100),
			Amount:        uint256.NewInt(60),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeBuy,
		},
		{
			Id:            2,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(100),
			Amount:        uint256.NewInt(40),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeBuy,
		},
	}
	asks := []*Order{
		{
			Id:            3,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(90),
			Amount:        uint256.NewInt(100),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeSell,
		},
	}
	orderBook := setupOrderBook(bids, asks)
	expectedTrades := []*Trade{{BidId: 1, AskId: 3}, {BidId: 2, AskId: 3}}
	trades, err := orderBook.MatchOrders()

	assert.NoError(t, err)
	assert.NotNil(t, trades)
	assert.Equal(t, expectedTrades, trades)
}

func TestBidNoMatchingAsk(t *testing.T) {
	bids := []*Order{
		{
			Id:            1,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(100),
			Amount:        uint256.NewInt(100),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeBuy,
		},
	}
	asks := []*Order{}
	orderBook := setupOrderBook(bids, asks)
	trades, err := orderBook.MatchOrders()

	assert.Error(t, err)
	assert.Nil(t, trades)
	assert.Equal(t, ErrNoMatch, err)
}

func TestAskNoMatchingBid(t *testing.T) {
	bids := []*Order{}
	asks := []*Order{
		{
			Id:            1,
			Hook:          testHook,
			SqrtPrice:     uint256.NewInt(100),
			Amount:        uint256.NewInt(100),
			MatchedAmount: uint256.NewInt(0),
			Type:          &OrderTypeSell,
		},
	}
	orderBook := setupOrderBook(bids, asks)
	trades, err := orderBook.MatchOrders()

	assert.Error(t, err)
	assert.Nil(t, trades)
	assert.Equal(t, ErrNoMatch, err)
}
