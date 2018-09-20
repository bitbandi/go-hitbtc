package main

import (
	"fmt"

	"github.com/saniales/go-hitbtc"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func main() {
	// hitbtc client
	hitbtc := hitbtc.New(API_KEY, API_SECRET)

	// GetBalances
	balances, _ := hitbtc.GetBalances()
	fmt.Println(len(balances))

	for _, bal := range balances {
		if bal.Currency == "BTC" {
			fmt.Println(bal.Available)
		}
	}

}
