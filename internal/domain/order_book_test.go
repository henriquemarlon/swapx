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
		{Id: 1, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(50), Type: &OrderTypeBuy},
	}
	asks := []*Order{
		{Id: 3, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(50), Type: &OrderTypeSell},
	}
	orderBook := setupOrderBook(bids, asks)

	expected := map[string][]string{"1": {"3"}}
	buyToSell, sellToBuy, err := orderBook.MatchOrders()
	assert.NoError(t, err)
	assert.Equal(t, expected, buyToSell)
	assert.Nil(t, sellToBuy)
}

func TestBidFullyMatchedByMultipleAsks(t *testing.T) {
	bids := []*Order{
		{Id: 1, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(100), Type: &OrderTypeBuy},
	}
	asks := []*Order{
		{Id: 3, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(40), Type: &OrderTypeSell},
		{Id: 4, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(60), Type: &OrderTypeSell},
	}
	orderBook := setupOrderBook(bids, asks)

	expected := map[string][]string{"1": {"3", "4"}}
	buyToSell, sellToBuy, err := orderBook.MatchOrders()
	assert.NoError(t, err)
	assert.Equal(t, expected, buyToSell)
	assert.Nil(t, sellToBuy)
}

func TestBidPartiallyMatched(t *testing.T) {
	bids := []*Order{
		{Id: 1, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(100), Type: &OrderTypeBuy},
	}
	asks := []*Order{
		{Id: 3, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(50), Type: &OrderTypeSell},
	}
	orderBook := setupOrderBook(bids, asks)

	expected := map[string][]string{"1": {"3"}}
	buyToSell, sellToBuy, err := orderBook.MatchOrders()
	assert.NoError(t, err)
	assert.Equal(t, expected, buyToSell)
	assert.Nil(t, sellToBuy)
}

func TestAskFullyMatchedBySingleBid(t *testing.T) {
	bids := []*Order{
		{Id: 1, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(50), Type: &OrderTypeBuy},
	}
	asks := []*Order{
		{Id: 2, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(50), Type: &OrderTypeSell},
	}
	orderBook := setupOrderBook(bids, asks)

	expected := map[string][]string{"2": {"1"}}
	buyToSell, sellToBuy, err := orderBook.MatchOrders()
	assert.NoError(t, err)
	assert.Nil(t, buyToSell)
	assert.Equal(t, expected, sellToBuy)
}

func TestAskFullyMatchedByMultipleBids(t *testing.T) {
	bids := []*Order{
		{Id: 1, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(40), Type: &OrderTypeBuy},
		{Id: 2, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(60), Type: &OrderTypeBuy},
	}
	asks := []*Order{
		{Id: 3, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(100), Type: &OrderTypeSell},
	}
	orderBook := setupOrderBook(bids, asks)

	expected := map[string][]string{"3": {"1", "2"}}
	buyToSell, sellToBuy, err := orderBook.MatchOrders()
	assert.NoError(t, err)
	assert.Nil(t, buyToSell)
	assert.Equal(t, expected, sellToBuy)
}

func TestBidNoMatchingAsk(t *testing.T) {
	bids := []*Order{
		{Id: 1, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(100), Type: &OrderTypeBuy},
	}
	asks := []*Order{}
	orderBook := setupOrderBook(bids, asks)

	buyToSell, sellToBuy, err := orderBook.MatchOrders()
	assert.Error(t, err)
	assert.Nil(t, buyToSell)
	assert.Nil(t, sellToBuy)
}

func TestAskNoMatchingBid(t *testing.T) {
	bids := []*Order{}
	asks := []*Order{
		{Id: 1, Hook: testHook, SqrtPrice: uint256.NewInt(100), Amount: uint256.NewInt(100), Type: &OrderTypeSell},
	}
	orderBook := setupOrderBook(bids, asks)

	buyToSell, sellToBuy, err := orderBook.MatchOrders()
	assert.Error(t, err)
	assert.Nil(t, buyToSell)
	assert.Nil(t, sellToBuy)
}
