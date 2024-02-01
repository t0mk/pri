package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/t0mk/pri/symbols"
)

const (
	Coinmate Exchange = "coinmate"
	Coinbase Exchange = "coinbase"
	Binance  Exchange = "binance"
	Kraken   Exchange = "kraken"
	Bitstamp Exchange = "bitstamp"
	Huobi    Exchange = "huobi"
	Kucoin   Exchange = "kucoin"
	Gateio   Exchange = "gateio"
	Bitfinex Exchange = "bitfinex"
	Bybit    Exchange = "bybit"
	Okx      Exchange = "okx"
)

var exchangeGetters = map[Exchange]TickerGetter{
	Coinmate: CoinmateGetter,
	Coinbase: CoinbaseGetter,
	Binance:  BinanceGetter,
	Kraken:   KrakenGetter,
	Bitstamp: BitstampGetter,
	Huobi:    HuobiGetter,
	Kucoin:   KUCoinGetter,
	Gateio:   GateIOGetter,
	Bitfinex: BitfinexGetter,
	Bybit:    BybitGetter,
	Okx:      OkxGetter,
}

var exchangeSymbols = map[Exchange][]string{
	Coinmate: symbols.Coinmate,
	Coinbase: symbols.Coinbase,
	Binance:  symbols.Binance,
	Kraken:   symbols.Kraken,
	Bitstamp: symbols.Bitstamp,
	Huobi:    symbols.Huobi,
	Kucoin:   symbols.Kucoin,
	Gateio:   symbols.Gateio,
	Bitfinex: symbols.Bitfinex,
	Bybit:    symbols.Bybit,
	Okx:      symbols.Okx,
}

var debug = false

type Exchange string
type ExTick struct {
	Exchange Exchange
	Ticker   string
}

type ExTickSet struct {
	Name    string
	ExTicks []ExTick
}

type TickPrice struct {
	MinAsk float64
	MaxBid float64
}

type ExTickPri struct {
	ExTick
	Price TickPrice
}

type ExTickPris []ExTickPri

func (etps ExTickPris) MinAsk() ExTickPri {
	min := etps[0]
	for _, v := range etps {
		if v.Price.MinAsk < min.Price.MinAsk {
			min = v
		}
	}
	return min
}

func (etps ExTickPris) MaxBid() ExTickPri {
	max := etps[0]
	for _, v := range etps {
		if v.Price.MaxBid > max.Price.MaxBid {
			max = v
		}
	}
	return max
}

func (etps ExTickPris) String() string {
	str := ""
	for _, v := range etps {
		str += fmt.Sprintf("%s\n", v)
	}
	return str
}

func ExTickPrisToArbResult(n string, etps ExTickPris) ArbResult {
	return ArbResult{
		Name:          n,
		TickPris:      etps,
		MinAsk:        etps.MinAsk(),
		MaxBid:        etps.MaxBid(),
		SpreadPercent: SpreadType(etps.SpreadPercent()),
	}
}

type SpreadType float64

func (st SpreadType) String() string {
	return fmt.Sprintf("%.3f%%", st)
}

type ArbResult struct {
	Name          string
	TickPris      ExTickPris
	MinAsk        ExTickPri
	MaxBid        ExTickPri
	SpreadPercent SpreadType
}

func (ar ArbResult) String() string {
	return fmt.Sprintf("name: %s\nmin: %s\nmax: %s\nspread: %s", ar.Name, ar.MinAsk, ar.MaxBid, ar.SpreadPercent)
}

type ArbResults []ArbResult

func (ars ArbResults) HighestSpread() ArbResult {
	max := ars[0]
	for _, v := range ars {
		if v.SpreadPercent > max.SpreadPercent {
			max = v
		}
	}
	return max
}

func (ars ArbResults) LowestSpread() ArbResult {
	min := ars[0]
	for _, v := range ars {
		if v.SpreadPercent < min.SpreadPercent {
			min = v
		}
	}
	return min
}

func (ars ArbResults) Report() {
	for _, v := range ars {
		fmt.Println(v.Name, v.SpreadPercent)
	}
	fmt.Println("Highest spread:")
	fmt.Println(ars.HighestSpread())
	fmt.Println("Lowest spread:")
	fmt.Println(ars.LowestSpread())
}

func (etps ExTickPris) SpreadPercent() float64 {
	min := etps.MinAsk().Price.MinAsk
	max := etps.MaxBid().Price.MaxBid
	return (max - min) / min * 100
}

func (etp ExTickPri) String() string {
	et := etp.ExTick
	return fmt.Sprintf("[%15s]\t[%15s - %15s]", et, formatFloat(etp.Price.MinAsk), formatFloat(etp.Price.MaxBid))
}

func (et ExTick) String() string {
	return fmt.Sprintf("%s-%s", et.Exchange, et.Ticker)
}

func getExchangeTickerPrice(et ExTick) (*ExTickPri, error) {
	start := time.Now()
	getter := exchangeGetters[et.Exchange]
	ab, err := getter(et.Ticker)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	if debug {
		log.Printf("[%15s]\t[%5s]\t%s\n", et, elapsed, ab)
	}
	return &ExTickPri{et, TickPrice{ab.Ask, ab.Bid}}, nil
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
		lst := []string{}
		for i, v := range foundSymbols {
			lst = append(lst, fmt.Sprintf("%s-%s", foundExchanges[i], v))
		}
		lines := strings.Join(lst, "\n")
		return nil, fmt.Errorf("symbol %s found in multiple exchanges:\n%s", symbol, lines)
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
	replacer := strings.NewReplacer("-", "", "/", "", "_", "")
	found := []ExTick{}
	for k, v := range exchangeSymbols {
		for _, t := range v {
			lowerT := strings.ToLower(t)
			lowerTNoDashNoSlash := replacer.Replace(lowerT)
			lowerSymbol := strings.ToLower(symbol)
			lowerSymbolNoDashNoSlash := replacer.Replace(lowerSymbol)
			if strings.Contains(lowerTNoDashNoSlash, lowerSymbolNoDashNoSlash) {
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
		if stringIsInSlice(os.Args[1], arbitrageCommands) {
			fmt.Println("arbitrage groups:")
			for k, v := range arbs {
				fmt.Printf("%s: %s\n", k, v)
			}
			return
		}
		if os.Args[1] == "aa" {
			exTickSets := []ExTickSet{}
			for name, a := range arbs {
				arbTicks, err := ExTicksFromSlice(a)
				if err != nil {
					panic(err)
				}
				if len(arbTicks) < 2 {
					fmt.Printf("Arbitrage group %s has less than 2 tickers, skipping\n", a)
					continue
				}
				exTickSets = append(exTickSets, ExTickSet{name, arbTicks})
			}
			arbResultChannel := getArbResultAsync(exTickSets)
			ars := CollectArbResultsFromChannel(arbResultChannel, len(exTickSets))

			ars.Report()
			return
		}

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
				tickerChannel := getExTickPriAsync(found)

				for _, v := range found {
					go getExchangeTickerPriceAsync(v, tickerChannel)
				}
				for i := 0; i < len(found); i++ {
					etp := <-tickerChannel
					fmt.Println(etp)
				}
			}
		} else if stringIsInSlice(cmd, arbitrageCommands) {
			arbName := os.Args[2]
			arb, ok := arbs[arbName]
			if !ok {
				fmt.Printf("Arbitrage group %s not found\n", arbName)
				return
			}
			for _, v := range arb {
				fmt.Println(v)
			}
			if strings.Contains(cmd, "!") {
				fmt.Println("Getting prices...")
				arbTicks, err := ExTicksFromSlice(arb)
				if err != nil {
					panic(err)
				}
				arbResult := ExTicksToArbResult(ExTickSet{arbName, arbTicks})
				fmt.Println(arbResult)

			}
		}
	}
}

func getArbResultAsync(exTickSets []ExTickSet) chan ArbResult {
	arbResultChannel := make(chan ArbResult)
	for _, v := range exTickSets {
		go func(ets ExTickSet) {
			arbResultChannel <- ExTicksToArbResult(ets)
		}(v)
	}
	return arbResultChannel
}

func CollectArbResultsFromChannel(channel chan ArbResult, n int) ArbResults {
	result := ArbResults{}
	for i := 0; i < n; i++ {
		ar := <-channel
		result = append(result, ar)
	}
	return result
}

func ExTicksToArbResult(ets ExTickSet) ArbResult {
	tickerChannel := getExTickPriAsync(ets.ExTicks)
	result := CollectExTickPrisFromChannel(tickerChannel, len(ets.ExTicks))
	arbResult := ExTickPrisToArbResult(ets.Name, result)
	return arbResult
}

func getExTickPriAsync(tickers []ExTick) chan *ExTickPri {
	tickerChannel := make(chan *ExTickPri)
	for _, v := range tickers {
		go getExchangeTickerPriceAsync(v, tickerChannel)
	}
	return tickerChannel
}

func CollectExTickPrisFromChannel(channel chan *ExTickPri, n int) ExTickPris {
	result := ExTickPris{}
	for i := 0; i < n; i++ {
		etp := <-channel
		fmt.Println(etp)
		result = append(result, *etp)
	}
	return result
}
