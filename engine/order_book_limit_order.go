package engine

import (
	"fmt"
	"order_engine/types"
)

// Process an order and return the trades generated before adding the remaining amount to the market
func (book *OrderBook) Process(order *types.Order) []types.Trade {
	if order.Side == "BUY" {
		return book.processLimitBuy(order)
	}
	return book.processLimitSell(order)
}

// Process a limit buy order
func (book *OrderBook) processLimitBuy(order *types.Order) []types.Trade {
	trades := make([]types.Trade, 0, 1)
	oskv := order.GetOBKeys()
	fmt.Println(oskv)
	pps, err := book.GetMatchingSellPricePoints(oskv)
	if err.Error() == "nil returned" {
		book.addOrder(order)
		return nil
	}

	if len(pps) == 0 {
		book.addOrder(order)
		return nil
	}
	for _, pp := range pps {
		fmt.Println("pps", pp)
	}
	// finally add the remaining order to the list
	book.addOrder(order)

	return trades
}

// Process a limit sell order
func (book *OrderBook) processLimitSell(order *types.Order) []types.Trade {
	trades := make([]types.Trade, 0, 1)
	obkv := order.GetOBKeys()
	fmt.Println(obkv)
	pps, err := book.GetMatchingBuyPricePoints(obkv)
	fmt.Println("orders", pps)
	if err.Error() == "nil returned" {
		book.addOrder(order)
		return nil
	}

	if len(pps) == 0 {
		book.addOrder(order)
		return nil
	}
	for _, pp := range pps {
		fmt.Println("pps", pp)
	}
	book.addOrder(order)
	return trades
}
