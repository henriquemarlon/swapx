package domain

import (
	"container/heap"
	"errors"
	"strconv"

	"github.com/holiman/uint256"
)

var (
	ErrNoMatch = errors.New("no match found")
)

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
	if h[i].SqrtPrice.Cmp(h[j].SqrtPrice) == 0 {
		return h[i].Amount.Cmp(h[j].Amount) > 0
	}
	return h[i].SqrtPrice.Cmp(h[j].SqrtPrice) > 0
}

func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(*Order)) }
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
	if h[i].SqrtPrice.Cmp(h[j].SqrtPrice) == 0 {
		return h[i].Amount.Cmp(h[j].Amount) > 0
	}
	return h[i].SqrtPrice.Cmp(h[j].SqrtPrice) < 0
}

func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(*Order)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old) - 1
	item := old[n]
	*h = old[:n]
	return item
}

func (ob *OrderBook) MatchOrders() (map[string][]string, map[string][]string, error) {
	buyToSell := make(map[string][]string)
	sellToBuy := make(map[string][]string)

	hasBuyMatch := false
	hasSellMatch := false

	for ob.Bids.Len() > 0 && ob.Asks.Len() > 0 {
		highestBid := heap.Pop(ob.Bids).(*Order)
		var matchedAsks []string

		for ob.Asks.Len() > 0 {
			lowestAsk := heap.Pop(ob.Asks).(*Order)

			if highestBid.SqrtPrice.Cmp(lowestAsk.SqrtPrice) < 0 {
				heap.Push(ob.Bids, highestBid)
				heap.Push(ob.Asks, lowestAsk)
				break
			}

			matchedAsks = append(matchedAsks, strconv.FormatUint(lowestAsk.Id, 10))
			hasBuyMatch = true

			matchSize := new(uint256.Int)
			if highestBid.Amount.Cmp(lowestAsk.Amount) < 0 {
				matchSize.Set(highestBid.Amount)
			} else {
				matchSize.Set(lowestAsk.Amount)
			}

			highestBid.Amount.Sub(highestBid.Amount, matchSize)
			lowestAsk.Amount.Sub(lowestAsk.Amount, matchSize)

			if highestBid.Amount.Sign() == 0 {
				break
			}

			if lowestAsk.Amount.Sign() > 0 {
				heap.Push(ob.Asks, lowestAsk)
			}
		}

		if len(matchedAsks) > 0 {
			buyToSell[strconv.FormatUint(highestBid.Id, 10)] = matchedAsks
		}

		if highestBid.Amount.Sign() > 0 {
			heap.Push(ob.Bids, highestBid)
		}
	}

	for ob.Asks.Len() > 0 && ob.Bids.Len() > 0 {
		lowestAsk := heap.Pop(ob.Asks).(*Order)
		var matchedBids []string

		for ob.Bids.Len() > 0 {
			highestBid := heap.Pop(ob.Bids).(*Order)

			if highestBid.SqrtPrice.Cmp(lowestAsk.SqrtPrice) < 0 {
				heap.Push(ob.Asks, lowestAsk)
				heap.Push(ob.Bids, highestBid)
				break
			}

			matchedBids = append(matchedBids, strconv.FormatUint(highestBid.Id, 10))
			hasSellMatch = true

			matchSize := new(uint256.Int)
			if lowestAsk.Amount.Cmp(highestBid.Amount) < 0 {
				matchSize.Set(lowestAsk.Amount)
			} else {
				matchSize.Set(highestBid.Amount)
			}

			lowestAsk.Amount.Sub(lowestAsk.Amount, matchSize)
			highestBid.Amount.Sub(highestBid.Amount, matchSize)

			if lowestAsk.Amount.Sign() == 0 {
				break
			}

			if highestBid.Amount.Sign() > 0 {
				heap.Push(ob.Bids, highestBid)
			}
		}

		if len(matchedBids) > 0 {
			sellToBuy[strconv.FormatUint(lowestAsk.Id, 10)] = matchedBids
		}
	}

	if !hasBuyMatch && !hasSellMatch {
		return nil, nil, ErrNoMatch
	}

	if hasBuyMatch && !hasSellMatch {
		return buyToSell, nil, nil
	}
	if hasSellMatch && !hasBuyMatch {
		return nil, sellToBuy, nil
	}

	return buyToSell, sellToBuy, nil
}
