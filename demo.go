package main

var btcUsdtTickers = []ExTick{
	{"coinbase", "BTC-USD"},
	{"binance", "BTCUSDT"},
	{"kraken", "XBTUSDT"},
	{"bitstamp", "btcusdt"},
	{"huobi", "btcusdt"},
	{"kucoin", "BTC-USDT"},
	{"gateio", "BTC_USDT"},
	{"bitfinex", "btcust"},
}

func floatSliceMean(slice []float64) float64 {
	sum := 0.0
	i := 0
	for _, v := range slice {
		if v == 0 {
			continue
		}
		sum += v
		i++
	}
	return sum / float64(i)
}

/*
func getBtcUsdt() float64 {
	prices := []float64{}
	for _, et := range btcUsdtTickers {
		price := getExchangeTickerPrice(et)
		prices = append(prices, price)
	}
	return floatSliceMean(prices)
}

func getBtcUsdtAsync() float64 {
	prices := []float64{}
	resultChan := make(chan float64)
	for _, et := range btcUsdtTickers {
		go getExchangeTickerPriceAsync(et, resultChan)
	}
	for i := 0; i < len(btcUsdtTickers); i++ {
		price := <-resultChan
		prices = append(prices, price)
	}
	return floatSliceMean(prices)
}
*/
