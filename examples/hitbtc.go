package main

import (
	"fmt"
	"github.com/sutapurachina/go-hitbtc"
)

const (
	ApiKey    = ""
	ApiSecret = ""
)

func main() {
	// hitbtc client
	hitbtc := hitbtc.New(ApiKey, ApiSecret)

	// GetBalances
	balances, _ := hitbtc.GetBalances()
	fmt.Println(len(balances))

	for _, bal := range balances {
		if bal.Currency == "BTC" {
			fmt.Println(bal.Available)
		}
	}

}
