package domain

import (
	"container/heap"
	"errors"

	"github.com/holiman/uint256"
)

var (
	ErrNoMatch       = errors.New("no match found")
	ErrOrderNotFound = errors.New("order not found")
)

type OrderBook struct {
	Bids *MaxHeap // Buy orders (highest price first)
	Asks *MinHeap // Sell orders (lowest price first)
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

// MaxHeap: Buy orders (highest price first)
type MaxHeap []*Order

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].SqrtPrice.Cmp(h[j].SqrtPrice) > 0 }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(*Order)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old) - 1
	item := old[n]
	*h = old[:n]
	return item
}

// MinHeap: Sell orders (lowest price first)
type MinHeap []*Order

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i].SqrtPrice.Cmp(h[j].SqrtPrice) < 0 }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(*Order)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old) - 1
	item := old[n]
	*h = old[:n]
	return item
}

func (ob *OrderBook) AddOrder(order *Order) {
	if order.Type == &OrderTypeBuy {
		heap.Push(ob.Bids, order)
	} else {
		heap.Push(ob.Asks, order)
	}
	ob.MatchOrders()
}

func (ob *OrderBook) RemoveOrder(order *Order) error {
	if order.Type == &OrderTypeBuy {
		for i, o := range *ob.Bids {
			if o.Id == order.Id {
				heap.Remove(ob.Bids, i)
				return nil
			}
		}
	} else {
		for i, o := range *ob.Asks {
			if o.Id == order.Id {
				heap.Remove(ob.Asks, i)
				return nil
			}
		}
	}
	return ErrOrderNotFound
}

func (ob *OrderBook) MatchOrders() {
	for ob.Bids.Len() > 0 && ob.Asks.Len() > 0 {
		highestBid := heap.Pop(ob.Bids).(*Order)
		lowestAsk := heap.Pop(ob.Asks).(*Order)

		if highestBid.SqrtPrice.Cmp(lowestAsk.SqrtPrice) < 0 {
			heap.Push(ob.Bids, highestBid)
			heap.Push(ob.Asks, lowestAsk)
			return // No match possible
		}

		matchSize := highestBid.Amount
		if lowestAsk.Amount.Cmp(highestBid.Amount) < 0 {
			matchSize = lowestAsk.Amount
		}

		highestBid.Amount = new(uint256.Int).Sub(highestBid.Amount, matchSize)
		lowestAsk.Amount = new(uint256.Int).Sub(lowestAsk.Amount, matchSize)

		if highestBid.Amount.Sign() > 0 {
			heap.Push(ob.Bids, highestBid)
		}
		if lowestAsk.Amount.Sign() > 0 {
			heap.Push(ob.Asks, lowestAsk)
		}
	}
}
