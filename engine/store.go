package engine

import "fmt"

// GetPricePoints returns the pricepoints matching a certain (pair, pricepoint)
func (ob *OrderBook) GetMatchingBuyPricePoints(obKey string) ([]string, error) {
	pps, err := ob.RedisConn.GetKeys(obKey)
	if err != nil {
		return nil, err
	}

	return pps, nil
}

func (ob *OrderBook) GetMatchingSellPricePoints(obkv string) ([]string, error) {
	pps, err := ob.RedisConn.GetKeys(obkv)
	if err != nil {
		return nil, err
	}

	return pps, nil
}

// AddPricePointToSet
func (ob *OrderBook) AddToPricePointSet(pricePointSetKey string, pricePoint uint64) error {
	fmt.Println(pricePointSetKey, pricePoint)
	err := ob.RedisConn.Set(pricePointSetKey, pricePoint)
	if err != nil {
		return err
	}

	return nil
}

func (ob *OrderBook) GetPricePoint(pricePointKey string) (uint64, error) {
	pp, err := ob.RedisConn.Get(pricePointKey)
	if err != nil {
		return 0, err
	}

	return pp, nil
}

func (ob *OrderBook) RemPricePoint(pricePointKey string) (error) {
	err := ob.RedisConn.Del(pricePointKey)
	if err != nil {
		return err
	}

	return nil
}
