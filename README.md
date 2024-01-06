# pri

CLI tool for simply getting cryptocurrency prices from 8 big exchanges. No API key needed.

Supported exchanges are:
- Binance
- Coinbase
- Kraken
- GateIO
- KuCoin
- BitStamp
- BitFinex
- Huobi

## Install

To install latest release for Linux:

```sh
wget -O pri https://github.com/t0mk/pri/releases/latest/download/pri-linux-amd64 && chmod +x pri && sudo cp pri /usr/local/bin/
```

.. for MacOS:

```sh
wget -O pri https://github.com/t0mk/pri/releases/latest/download/pri-darwin-amd64 && chmod +x pri && sudo cp pri /usr/local/bin/
```

## Build

```sh
git clone https://github.com/t0mk/pri
cd pri
go build
```

## Usage

Get price of a crypto pair, for example BTC-USDT

```zsh
➜  ./pri BTCUSDT
[          binance-BTCUSDT]	43,952.19
➜  ./pri btcusdt
[            huobi-btcusdt]	43,937.54
➜  ./pri BTCUSDT
[          binance-BTCUSDT]	43,934.00
➜  ./pri kraken-XBTUSDT 
[           kraken-XBTUSDT]	43,919.40
```

If you pass the pair without exchange prefix, `pri` will try to find it in known exchange assets.

You can also lookup all crypto pairs containing "rpl":

```zsh
➜  p./pri find rpl
Found symbols:
coinbase-RPL-USD
binance-RPLBTC
binance-RPLBUSD
binance-RPLUSDT
kraken-RPLEUR
kraken-RPLUSD
huobi-rplusdt
kucoin-RPL-USDT
gateio-RPL_USDT
```

If you add `!` after the find command, `pri` will fetch prices for the found assets:

```zsh
➜ ./pri find! rpl
Found symbols:
kucoin-RPL-USDT
gateio-RPL_USDT
coinbase-RPL-USD
binance-RPLBTC
binance-RPLBUSD
binance-RPLUSDT
kraken-RPLEUR
kraken-RPLUSD
huobi-rplusdt
Getting prices...
[         coinbase-RPL-USD]	29.97
[            kraken-RPLEUR]	27.46
[            kraken-RPLUSD]	29.99
[            huobi-rplusdt]	29.90
[          kucoin-RPL-USDT]	29.95
[          gateio-RPL_USDT]	29.96
[          binance-RPLBUSD]	24.17
[           binance-RPLBTC]	0.000641
[          binance-RPLUSDT]	29.94
```


## TODO
- add bybit
- automate asset list compilation
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
