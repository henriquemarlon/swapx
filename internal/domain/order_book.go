package domain

import (
	"container/heap"
	"errors"

	"github.com/holiman/uint256"
)

var ErrNoMatch = errors.New("no match found")

type Trade struct {
	BidId uint64
	AskId uint64
}

type OrderBook struct {
	Bids *MaxHeap
	Asks *MinHeap
}

func NewOrderBook() *OrderBook {
	bids := &MaxHeap{}
	asks := &MinHeap{}
	heap.Init(bids)
	heap.Init(asks)
	return &OrderBook{
		Bids: bids,
		Asks: asks,
	}
}

type MaxHeap []*Order

func (h MaxHeap) Len() int { return len(h) }

func (h MaxHeap) Less(i, j int) bool {
	priceCmp := h[i].SqrtPrice.Cmp(h[j].SqrtPrice)
	if priceCmp != 0 {
		return priceCmp > 0
	}
	return h[i].Id < h[j].Id
}

func (h MaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(*Order))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old) - 1
	item := old[n]
	*h = old[:n]
	return item
}

type MinHeap []*Order

func (h MinHeap) Len() int { return len(h) }

func (h MinHeap) Less(i, j int) bool {
	priceCmp := h[i].SqrtPrice.Cmp(h[j].SqrtPrice)
	if priceCmp != 0 {
		return priceCmp < 0
	}
	return h[i].Id < h[j].Id
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(*Order))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old) - 1
	item := old[n]
	*h = old[:n]
	return item
}

func (ob *OrderBook) MatchOrders() ([]*Trade, error) {
	var trades []*Trade

	for ob.Bids.Len() > 0 && ob.Asks.Len() > 0 {
		bestBid := (*ob.Bids)[0]
		bestAsk := (*ob.Asks)[0]
	
		if bestBid.SqrtPrice.Cmp(bestAsk.SqrtPrice) < 0 {
			break
		}
	
		remainingBid := new(uint256.Int).Sub(bestBid.Amount, bestBid.MatchedAmount)
		remainingAsk := new(uint256.Int).Sub(bestAsk.Amount, bestAsk.MatchedAmount)
	
		matchedQty := new(uint256.Int)
		if remainingBid.Cmp(remainingAsk) <= 0 {
			matchedQty.Set(remainingBid)
		} else {
			matchedQty.Set(remainingAsk)
		}
			
		trade := &Trade{
			BidId: bestBid.Id,
			AskId: bestAsk.Id,
		}
		trades = append(trades, trade)
	
		bestBid.MatchedAmount = new(uint256.Int).Add(bestBid.MatchedAmount, matchedQty)
		bestAsk.MatchedAmount = new(uint256.Int).Add(bestAsk.MatchedAmount, matchedQty)
	
		remainingBid = new(uint256.Int).Sub(bestBid.Amount, bestBid.MatchedAmount)
		remainingAsk = new(uint256.Int).Sub(bestAsk.Amount, bestAsk.MatchedAmount)
	
		if remainingBid.IsZero() {
			heap.Pop(ob.Bids)
		}
	
		if remainingAsk.IsZero() {
			heap.Pop(ob.Asks)
		}
	}

	if len(trades) == 0 {
		return nil, ErrNoMatch
	}
	return trades, nil
}