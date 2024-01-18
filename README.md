# pri

CLI tool for simply getting cryptocurrency prices from biggest exchanges. No API key needed.

Supported exchanges are:
- Binance
- Coinbase
- Kraken
- GateIO
- KuCoin
- BitStamp
- BitFinex
- Huobi
- Bybit

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

You can also check arbitrages - same ticker traded on different exchanges. `pri` will report highest and lowest price and spread.
```
./pri a! btc-usdt
coinbase-BTC-USDT
huobi-btcusdt
kucoin-BTC-USDT
binance-BTCUSDT
gateio-BTC_USDT
bitfinex-BTCUST
bybit-BTCUSDT
kraken-XBTUSDT
okx-BTC-USDT
Getting prices...
[bitfinex-BTCUST]	41,893.00
[coinbase-BTC-USDT]	41,883.71
[  bybit-BTCUSDT]	41,894.25
[ kraken-XBTUSDT]	41,889.45
[   okx-BTC-USDT]	41,887.95
[  huobi-btcusdt]	41,894.04
[kucoin-BTC-USDT]	41,891.95
[gateio-BTC_USDT]	41,892.05
[binance-BTCUSDT]	41,892.10
name: btc-usdt
min: [coinbase-BTC-USDT]	41,883.71
max: [  bybit-BTCUSDT]	41,894.25
spread: 0.025%
```




## TODO
- automate asset list compilation
