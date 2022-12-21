package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func get() map[string]*HighLow {
	// Jan 1993 - Dec 2022
	url := "https://query1.finance.yahoo.com/v7/finance/download/SPY?period1=724889583&period2=1671574383&interval=1d&events=history"
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	bytes, _ := io.ReadAll(res.Body)
	body := string(bytes)
	highLow := make(map[string]*HighLow)
	for _, val := range strings.Split(body, "\n")[1:] {
		items := strings.Split(val, ",")
		dates := strings.Split(items[0], "-")
		key := fmt.Sprintf("%s-%s", dates[0], dates[1])

		high, _ := strconv.ParseFloat(items[2], 64)
		low, _ := strconv.ParseFloat(items[3], 64)

		val, ok := highLow[key]
		if ok {
			val.High = math.Max(val.High, high)
			val.Low = math.Min(val.Low, low)
		} else {
			highLow[key] = &HighLow{high, low}
		}
	}
	return highLow
}

func printStruct(s map[string]*HighLow) {
	for k, v := range s {
		fmt.Printf("%s: high %.2f, low %.2f\n", k, v.High, v.Low)
	}
}

func main() {
	highLow := get()
	high := 0.0
	low := 0.0
	counter := 0
	monthly := 200.0

	// SPY as of 12/21/2022
	curr_price := 386.10

	for year := 1993; year < 2023; year++ {
		for month := 1; month <= 12; month++ {
			key := fmt.Sprintf("%d-%02d", year, month)
			val, ok := highLow[key]
			if ok {
				fmt.Printf("%s, high: %.2f, low: %.2f\n", key, val.High, val.Low)
				high += monthly / val.High
				low += monthly / val.Low
			} else {
				panic("Key not found")
			}
			counter++
		}
	}
	fmt.Printf("Timing High: %.2f\n", high*curr_price)
	fmt.Printf("Timing Low: %.2f\n", low*curr_price)
	fmt.Printf("Cash: %d * %d months = %d\n", 200, counter, 200*counter)
}
