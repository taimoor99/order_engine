package engine

// GetPricePoints returns the pricepoints matching a certain (pair, pricepoint)
func (ob *OrderBook) GetMatchingBuyPricePoints(obKey string) ([]string, error) {
	pps, err := ob.RedisConn.Get(obKey)
	if err != nil {
		return nil, err
	}

	return pps, nil
}

func (ob *OrderBook) GetMatchingSellPricePoints(obkv string) ([]string, error) {
	pps, err := ob.RedisConn.Get(obkv)
	if err != nil {
		return nil, err
	}

	return pps, nil
}

// AddPricePointToSet
func (ob *OrderBook) AddToPricePointSet(pricePointSetKey string, pricePoint string) error {
	err := ob.RedisConn.Set(pricePointSetKey, pricePoint)
	if err != nil {
		return err
	}

	return nil
}
