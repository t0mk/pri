# pri

CLI tool for getting cryptocurrency prices from 8 biggest exchanges. No API key needed


## TODO
- filter out delisted asset pairs, it's sometimes indicated in the api asset listsing
- make arbitrage groups, like 
```golang
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
```
