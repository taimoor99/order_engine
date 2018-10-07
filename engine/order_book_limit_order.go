package engine

import (
	"order_engine/types"
)

// Process an order and return the trades generated before adding the remaining amount to the market
func (book *OrderBook) Process(order *types.Order) []types.Trade {
	if order.Side == "BUY" {
		return book.processBuy(order)
	}
	return book.processSell(order)
}

// Process a buy order
func (book *OrderBook) processBuy(order *types.Order) []types.Trade {
	trades := make([]types.Trade, 0, 1)
	oskv := order.GetOBKeys()
	pps, err := book.GetMatchingSellPricePoints(oskv)
	if err != nil {
		return nil
	}

	if len(pps) == 0 {
		book.addOrder(order)
		return nil
	}
	for _, ppk := range pps {
		pp, err := book.GetPricePoint(ppk)
		if err != nil {
			return nil
		}
		if pp == order.Price {
			if err := book.RemPricePoint(ppk); err != nil {
				return nil
			}
			trades = append(trades, types.Trade{
				MakerOrderID: order.TokenID,
				Price:        order.Price,
			})
			return trades
		}
	}
	// finally add the remaining order to the list
	book.addOrder(order)

	return trades
}

// Process a sell order
func (book *OrderBook) processSell(order *types.Order) []types.Trade {
	trades := make([]types.Trade, 0, 1)
	obkv := order.GetOBKeys()
	pps, err := book.GetMatchingBuyPricePoints(obkv)
	if err != nil {
		return nil
	}

	if len(pps) == 0 {
		book.addOrder(order)
		return nil
	}
	for _, ppk := range pps {
		pp, err := book.GetPricePoint(ppk)
		if err != nil {
			return nil
		}
		if pp == order.Price {
			if err := book.RemPricePoint(ppk); err != nil {
				return nil
			}
			trades = append(trades, types.Trade{
				MakerOrderID: order.TokenID,
				Price:        order.Price,
			})
			return trades
		}
	}
	book.addOrder(order)
	return trades
}
