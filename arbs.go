package main

import (
	"fmt"
	"strings"
)

func ExTicksFromSlice(ticks []string) ([]ExTick, error) {
	var exTicks []ExTick
	for _, tick := range ticks {
		split := strings.SplitN(tick, "-", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid tick: %s", tick)
		}
		exchange := split[0]
		ticker := split[1]
		exTick := ExTick{Exchange(exchange), ticker}
		exTicks = append(exTicks, exTick)
	}
	return exTicks, nil
}


var arbs = map[string][]string{
	"eth-btc": {
		"binance-ETHBTC",
		"bitfinex-ETHBTC",
		"coinbase-ETH-BTC",
		"gateio-ETH_BTC",
		"huobi-ethbtc",
		"kraken-XETHXXBT",
		"kucoin-ETH-BTC",
		"okx-ETH-BTC",
		//"bitstamp-ETH/BTC",
	},
	"btc-usdt": {
		"coinbase-BTC-USDT",
		//"bitstamp-BTC/USDT",
		"huobi-btcusdt",
		"kucoin-BTC-USDT",
		"binance-BTCUSDT",
		"gateio-BTC_USDT",
		"bitfinex-BTCUST",
		"bybit-BTCUSDT",
		"kraken-XBTUSDT",
		"okx-BTC-USDT",
	},
	"eth-usdt": {
		"binance-ETHUSDT",
		"bitfinex-ETHUST",
		// bitstamp-ETH/USDT,
		"bybit-ETHUSDT",
		"coinbase-ETH-USDT",
		"gateio-ETH_USDT",
		"huobi-ethusdt",
		"kraken-ETHUSDT",
		"kucoin-ETH-USDT",
		"okx-ETH-USDT",
	},
	"bnb-btc": {
		"binance-BNBBTC",
		"kucoin-BNB-BTC",
		"gateio-BNB_BTC",
	},
	"bnb-usdt": {
		"kucoin-BNB-USDT",
		"bybit-BNBUSDT",
		"okx-BNB-USDT",
		"binance-BNBUSDT",
		"huobi-bnbusdt",
		"gateio-BNB_USDT",
	},
	"sol-btc": {
		"coinbase-SOL-BTC",
		"huobi-solbtc",
		"binance-SOLBTC",
		"bitfinex-SOLBTC",
		"okx-SOL-BTC",
	},
	"sol-usdt": {
		"kraken-SOLUSDT",
		"kucoin-SOL-USDT",
		"gateio-SOL_USDT",
		"okx-SOL-USDT",
		"bitfinex-SOLUST",
		"bybit-SOLUSDT",
		"coinbase-SOL-USDT",
		"binance-SOLUSDT",
		"huobi-solusdt",
	},
	"sol-eth": {
		"coinbase-SOL-ETH",
		"kraken-SOLETH",
		"okx-SOL-ETH",
		"binance-SOLETH",
	},
	"xrp-btc": {
		//"bitstamp-XRP/BTC",
		"huobi-xrpbtc",
		"kucoin-XRP-BTC",
		"gateio-XRP_BTC",
		"bitfinex-XRPBTC",
		"binance-XRPBTC",
		"okx-XRP-BTC",
	},
	"xrp-usdt": {
		"binance-XRPUSDT",
		"huobi-xrpusdt",
		"kucoin-XRP-USDT",
		"okx-XRP-USDT",
		"coinbase-XRP-USDT",
		//"bitstamp-XRP/USDT",
		"gateio-XRP_USDT",
		"bitfinex-XRPUST",
		"bybit-XRPUSDT",
		"kraken-XRPUSDT",
	},
	"btc-usdc":  {},
	"ada-btc":   {},
	"ada-usdt":  {},
	"avax-btc":  {},
	"avax-usdt": {},
	"doge-btc":  {},
	"doge-usdt": {},
	"trx-btc":   {},
	"trx-usdt":  {},
	"dot-btc":   {},
	"dot-usdt":  {},
}
