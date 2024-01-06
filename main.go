package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/t0mk/pri/symbols"
)

var debug = false

type Exchange string
type ExTick struct {
	Exchange Exchange
	Ticker   string
}

type ExTickPri struct {
	ExTick
	Price float64
}

func (etp ExTickPri) String() string {
	et := etp.ExTick
	return fmt.Sprintf("[%25s]\t%s", et, formatFloat(etp.Price))
}

func (et ExTick) String() string {
	return fmt.Sprintf("%s-%s", et.Exchange, et.Ticker)
}

const (
	Coinbase Exchange = "coinbase"
	Binance  Exchange = "binance"
	Kraken   Exchange = "kraken"
	Bitstamp Exchange = "bitstamp"
	Huobi    Exchange = "huobi"
	Kucoin   Exchange = "kucoin"
	Gateio   Exchange = "gateio"
	Bitfinex Exchange = "bitfinex"
)

var exchangeGetters = map[Exchange]TickerGetter{
	Coinbase: CoinbaseGetter,
	Binance:  BinanceGetter,
	Kraken:   KrakenGetter,
	Bitstamp: BitstampGetter,
	Huobi:    HuobiGetter,
	Kucoin:   KUCoinGetter,
	Gateio:   GateIOGetter,
	Bitfinex: BitfinexGetter,
}

func getExchangeTickerPrice(et ExTick) (*ExTickPri, error) {
	start := time.Now()
	getter := exchangeGetters[et.Exchange]
	price, err := getter(et.Ticker)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	if debug {
		log.Printf("[%15s]\t[%5s]\t%.2f\n", et, elapsed, price)
	}
	return &ExTickPri{et, price}, nil
}

func getExchangeTickerPriceAsync(et ExTick, channel chan *ExTickPri) {
	etp, err := getExchangeTickerPrice(et)
	if err != nil {
		log.Printf("error getting price for %s: %v", et, err)
		channel <- nil
		return
	}
	channel <- etp
}

var exchangeSymbols = map[Exchange][]string{
	Coinbase: symbols.Coinbase,
	Binance:  symbols.Binance,
	Kraken:   symbols.Kraken,
	Bitstamp: symbols.Bitstamp,
	Huobi:    symbols.Huobi,
	Kucoin:   symbols.Kucoin,
	Gateio:   symbols.Gateio,
	Bitfinex: symbols.Bitfinex,
}

func findExTick(symbol string) (*ExTick, error) {
	if hasExchangePrefix(symbol) {
		sli := strings.SplitN(symbol, "-", 2)
		if len(sli) != 2 {
			return nil, fmt.Errorf("invalid symbol: %s", symbol)
		}
		exchange := Exchange(sli[0])
		ticker := sli[1]
		if _, ok := exchangeGetters[exchange]; !ok {
			return nil, fmt.Errorf("ticker %s not found in exchange %s", ticker, exchange)
		}
		return &ExTick{exchange, ticker}, nil
	}
	foundSymbols := []string{}
	foundExchanges := []Exchange{}
	for k, v := range exchangeSymbols {
		for _, t := range v {
			if t == symbol {
				foundSymbols = append(foundSymbols, t)
				foundExchanges = append(foundExchanges, k)
			}
		}
	}
	if len(foundSymbols) == 0 {
		return nil, fmt.Errorf("symbol not found: %s", symbol)
	}
	if len(foundSymbols) > 1 {
		return nil, fmt.Errorf("symbol \"%s\" found in multiple exchanges: %v ", symbol, foundExchanges)
	}
	return &ExTick{foundExchanges[0], foundSymbols[0]}, nil
}

func preSetup() {
	if (os.Getenv("DEBUG") != "") && (os.Getenv("DEBUG") != "0") {
		debug = true
	}
}

func stringIsInSlice(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

func searchForExchangeTicker(symbol string) []ExTick {
	found := []ExTick{}
	for k, v := range exchangeSymbols {
		for _, t := range v {
			lowerT := strings.ToLower(t)
			lowerSymbol := strings.ToLower(symbol)
			if strings.Contains(lowerT, lowerSymbol) {
				found = append(found, ExTick{k, t})
			}
		}
	}
	return found
}

func hasExchangePrefix(symbol string) bool {
	for e, _ := range exchangeSymbols {
		if strings.HasPrefix(symbol, string(e)) {
			return true
		}
	}
	return false
}

func main() {
	preSetup()
	if len(os.Args) == 1 {
		panic("usage: pri [find] <symbol>")
	}
	if len(os.Args) == 2 {
		xt, err := findExTick(os.Args[1])
		if err != nil {
			panic(err)
		}
		etp, err := getExchangeTickerPrice(*xt)
		if err != nil {
			panic(err)
		}
		fmt.Println(etp)
	}
	if len(os.Args) == 3 {
		cmd := os.Args[1]
		if stringIsInSlice(cmd, findCommands) {
			symbol := os.Args[2]
			if hasExchangePrefix(symbol) {
				sli := strings.SplitN(symbol, "-", 2)
				fmt.Printf("Symbol should not have exchange prefix, use only %s instead of %s\n", sli[1], symbol)
				return
			}
			if len(symbol) < 3 {
				fmt.Println("Symbol should be at least 3 characters long")
			}
			found := searchForExchangeTicker(symbol)
			if len(found) == 0 {
				fmt.Printf("Symbol \"%s\" not found\n", symbol)
				return
			}
			fmt.Println("Found symbols:")
			for _, v := range found {
				fmt.Println(v)
			}
			if strings.Contains(cmd, "!") {
				fmt.Println("Getting prices...")
				tickerChannel := make(chan *ExTickPri)
				for _, v := range found {
					go getExchangeTickerPriceAsync(v, tickerChannel)
				}
				for i := 0; i < len(found); i++ {
					etp := <-tickerChannel
					fmt.Println(etp)
				}
			}

		}
	}
}
