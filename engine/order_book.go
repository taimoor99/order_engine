package engine

import "fmt"

// OrderBook type
type OrderBook struct {
	BuyOrders  []Order
	SellOrders []Order
}

// Add a buy order to the order book
func (book *OrderBook) addBuyOrder(order Order) {
	n := len(book.BuyOrders)
	if n == 0 {
		book.BuyOrders = append(book.BuyOrders, order)
	}
	for i := n - 1; i >= 0; i-- {
		buyOrder := book.BuyOrders[i]
		if buyOrder.Price < order.Price {
			break
		}

		if i == n-1 {
			book.BuyOrders = append(book.BuyOrders, order)
		} else {
			copy(book.BuyOrders[i+1:], book.BuyOrders[i:])
			book.BuyOrders[i] = order
		}
	}
	fmt.Println("buy order book", len(book.BuyOrders))
}

// Add a sell order to the order book
func (book *OrderBook) addSellOrder(order Order) {
	n := len(book.SellOrders)
	if n == 0 {
		book.SellOrders = append(book.SellOrders, order)
	}
	for i := n - 1; i >= 0; i-- {
		sellOrder := book.SellOrders[i]
		if sellOrder.Price > order.Price {
			break
		}
		if i == n-1 {
			book.SellOrders = append(book.SellOrders, order)
			fmt.Println(i, n-1)
		} else {
			copy(book.SellOrders[i+1:], book.SellOrders[i:])
			book.SellOrders[i] = order
		}
	}
	fmt.Println("sell order book", len(book.SellOrders))
}

// Remove a buy order from the order book at a given index
func (book *OrderBook) removeBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

// Remove a sell order from the order book at a given index
func (book *OrderBook) removeSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
