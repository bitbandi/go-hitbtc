package main

import (
	"fmt"

	"github.com/bitbandi/go-hitbtc"
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

	for i, _ := range balances {
		if balances[i].Currency == "BTC" {
			fmt.Println(balances[i].Value)
		}
	}

}
