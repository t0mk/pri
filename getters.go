package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type TickerGetter func(string) (float64, error)

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

type CoinbaseTicker struct {
	Price  string `json:"price"`
	Ask    string `json:"ask"`
	Bid    string `json:"bid"`
	Volume string `json:"volume"`
	Size   string `json:"size"`
}

// ticker might be BTC-USD{T}
func CoinbaseGetter(ticker string) (float64, error) {
	url := "https://api.pro.coinbase.com/products/" + ticker + "/ticker"
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, err
	}
	var tickerData CoinbaseTicker
	err = json.Unmarshal(body, &tickerData)
	//fmt.Printf("Coinbase %s %#v\n", ticker, tickerData)
	if err != nil {
		return 0, err
	}
	a := tickerData.Ask
	b := tickerData.Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, err
	}
	return (af + bf) / 2, nil
}

func GeminiGetter(ticker string) (float64, error) {
	url := "https://api.gemini.com/v1/pubticker/" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, err
	}
	var tickerData CoinbaseTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return 0, err
	}
	a := tickerData.Ask
	b := tickerData.Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, err
	}
	return (af + bf) / 2, nil
}

type BitstampTicker struct {
	Ask string `json:"ask"`
	Bid string `json:"bid"`
}

func BitstampGetter(ticker string) (float64, error) {
	url := "https://www.bitstamp.net/api/v2/ticker/" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, err
	}
	var tickerData BitstampTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return 0, err
	}
	a := tickerData.Ask
	b := tickerData.Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, err
	}
	return (af + bf) / 2, nil
}

type KrakenTicker struct {
	Result map[string]struct {
		A []string `json:"a"`
		B []string `json:"b"`
	} `json:"result"`
}

func KrakenGetter(ticker string) (float64, error) {
	url := "https://api.kraken.com/0/public/Ticker?pair=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, err
	}
	var tickerData KrakenTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return 0, err
	}
	a := tickerData.Result[ticker].A[0]
	b := tickerData.Result[ticker].B[0]
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, err
	}
	return (af + bf) / 2, nil
}

type HuobiTicker struct {
	Tick struct {
		Bid []float64 `json:"bid"`
		Ask []float64 `json:"ask"`
	} `json:"tick"`
}

func HuobiGetter(ticker string) (float64, error) {
	url := "https://api.huobi.pro/market/detail/merged?symbol=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, err
	}
	var tickerData HuobiTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return 0, err
	}
	a := tickerData.Tick.Ask[0]
	b := tickerData.Tick.Bid[0]
	return (a + b) / 2, nil
}

func BinanceGetter(ticker string) (float64, error) {
	url := "https://api1.binance.com/api/v3/ticker/price?symbol=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, err
	}
	var tickerData CoinbaseTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return 0, err
	}
	price, err := strconv.ParseFloat(tickerData.Price, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

type GateIOTicker struct {
	Bid string `json:"highest_bid"`
	Ask string `json:"lowest_ask"`
}

func GateIOGetter(ticker string) (float64, error) {
	url := "https://api.gateio.ws/api/v4/spot/tickers?currency_pair=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, err
	}
	var tickerData []GateIOTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return 0, err
	}
	a := tickerData[0].Ask
	b := tickerData[0].Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, err
	}
	return (af + bf) / 2, nil
}

type BitfinexTicker struct {
	Mid string `json:"mid"`
}

func BitfinexGetter(ticker string) (float64, error) {
	url := "https://api.bitfinex.com/v1/pubticker/" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, err
	}
	var tickerData BitfinexTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return 0, err
	}
	price, err := strconv.ParseFloat(tickerData.Mid, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

type KUCoinTicker struct {
	Data struct {
		Price string `json:"price"`
		Bid   string `json:"bestBid"`
		Ask   string `json:"bestAsk"`
	} `json:"data"`
}

func KUCoinGetter(ticker string) (float64, error) {
	url := "https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, err
	}
	var tickerData KUCoinTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return 0, err
	}
	a := tickerData.Data.Ask
	b := tickerData.Data.Bid
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, err
	}
	return (af + bf) / 2, nil
}

type BitflyerTicker struct {
	Ask float64 `json:"best_ask"`
	Bid float64 `json:"best_bid"`
}

func BitflyerGetter(ticker string) (float64, error) {
	url := "https://api.bitflyer.com/v1/ticker?product_code=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, err
	}
	var tickerData BitflyerTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return 0, err
	}
	return (tickerData.Ask + tickerData.Bid) / 2, nil
}
