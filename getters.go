package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type AskBid struct {
	Ask float64
	Bid float64
}

type TickerGetter func(string) (*AskBid, error)

func getHTTPResponseBodyFromUrl(url string) ([]byte, error) {
	if debug {
		log.Println("getHTTPResponseBodyFromUrl", url)
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	if debug {
		log.Println("BODY:", string(body))
	}
	return body, nil
}

type CoinmateTicker struct {
	Data struct {
		Ask float64 `json:"ask"`
		Bid float64 `json:"bid"`
	} `json:"data"`
}

func CoinmateGetter(ticker string) (*AskBid, error) {
	url := "https://coinmate.io/api/ticker?currencyPair=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData CoinmateTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	return &AskBid{tickerData.Data.Ask, tickerData.Data.Bid}, nil
}

type CoinbaseTicker struct {
	Price  string `json:"price"`
	Ask    string `json:"ask"`
	Bid    string `json:"bid"`
	Volume string `json:"volume"`
	Size   string `json:"size"`
}

// ticker might be BTC-USD{T}
func CoinbaseGetter(ticker string) (*AskBid, error) {
	url := "https://api.pro.coinbase.com/products/" + ticker + "/ticker"
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData CoinbaseTicker
	err = json.Unmarshal(body, &tickerData)
	//fmt.Printf("Coinbase %s %#v\n", ticker, tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData.Ask
	b := tickerData.Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{af, bf}, nil
}

func GeminiGetter(ticker string) (*AskBid, error) {
	url := "https://api.gemini.com/v1/pubticker/" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData CoinbaseTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData.Ask
	b := tickerData.Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{af, bf}, nil
}

type BybitTicker struct {
	Reusult struct {
		List []struct {
			Ask string `json:"ask1Price"`
			Bid string `json:"bid1Price"`
		} `json:"list"`
	} `json:"result"`
}

func BybitGetter(ticker string) (*AskBid, error) {
	url := "https://api-testnet.bybit.com/v5/market/tickers?category=linear&symbol=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData BybitTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData.Reusult.List[0].Ask
	b := tickerData.Reusult.List[0].Bid

	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{af, bf}, nil
}

type BitstampTicker struct {
	Ask string `json:"ask"`
	Bid string `json:"bid"`
}

func BitstampGetter(ticker string) (*AskBid, error) {
	url := "https://www.bitstamp.net/api/v2/ticker/" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Bitstamp %s %s\n", ticker, string(body))
	//fmt.Println(url)
	var tickerData BitstampTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData.Ask
	b := tickerData.Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{af, bf}, nil
}

type KrakenTicker struct {
	Result map[string]struct {
		A []string `json:"a"`
		B []string `json:"b"`
	} `json:"result"`
}

func KrakenGetter(ticker string) (*AskBid, error) {
	url := "https://api.kraken.com/0/public/Ticker?pair=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData KrakenTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData.Result[ticker].A[0]
	b := tickerData.Result[ticker].B[0]
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{af, bf}, nil
}

type HuobiTicker struct {
	Tick struct {
		Bid []float64 `json:"bid"`
		Ask []float64 `json:"ask"`
	} `json:"tick"`
}

func HuobiGetter(ticker string) (*AskBid, error) {
	url := "https://api.huobi.pro/market/detail/merged?symbol=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData HuobiTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData.Tick.Ask[0]
	b := tickerData.Tick.Bid[0]
	return &AskBid{a, b}, nil
}

func BinanceGetter(ticker string) (*AskBid, error) {
	url := "https://api1.binance.com/api/v3/ticker/price?symbol=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData CoinbaseTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	price, err := strconv.ParseFloat(tickerData.Price, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{price, price}, nil
}

type GateIOTicker struct {
	Bid string `json:"highest_bid"`
	Ask string `json:"lowest_ask"`
}

func GateIOGetter(ticker string) (*AskBid, error) {
	url := "https://api.gateio.ws/api/v4/spot/tickers?currency_pair=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData []GateIOTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData[0].Ask
	b := tickerData[0].Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{af, bf}, nil
}

type BitfinexTicker struct {
	Mid string `json:"mid"`
}

func BitfinexGetter(ticker string) (*AskBid, error) {
	url := "https://api.bitfinex.com/v1/pubticker/" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData BitfinexTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	price, err := strconv.ParseFloat(tickerData.Mid, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{price, price}, nil
}

type KUCoinTicker struct {
	Data struct {
		Price string `json:"price"`
		Bid   string `json:"bestBid"`
		Ask   string `json:"bestAsk"`
	} `json:"data"`
}

func KUCoinGetter(ticker string) (*AskBid, error) {
	url := "https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData KUCoinTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData.Data.Ask
	b := tickerData.Data.Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{af, bf}, nil
}

type BitflyerTicker struct {
	Ask float64 `json:"best_ask"`
	Bid float64 `json:"best_bid"`
}

func BitflyerGetter(ticker string) (*AskBid, error) {
	url := "https://api.bitflyer.com/v1/ticker?product_code=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData BitflyerTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	return &AskBid{tickerData.Ask, tickerData.Bid}, nil
}

type OkxTicker struct {
	Data []struct {
		Ask string `json:"askPx"`
		Bid string `json:"bidPx"`
	} `json:"data"`
}

func OkxGetter(ticker string) (*AskBid, error) {
	url := "https://www.okx.com/api/v5/market/ticker?instId=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData OkxTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData.Data[0].Ask
	b := tickerData.Data[0].Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{af, bf}, nil
}
