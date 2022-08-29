go-hitbtc
==========

go-hitbtc is an implementation of the HitBTC API (public and private) in Golang.

This version implement V2 HitBTC API.

## Import
	go get "github.com/sutapurachina/go-hitbtc"
	
## Usage

In order to use the client with go's default http client settings you can do:

~~~ go
package main

import (
	"fmt"
	"github.com/sutapurachina/go-hitbtc"
)

const (
	API_KEY    = "YOUR_API_KEY"
	API_SECRET = "YOUR_API_SECRET"
)

func main() {
	// hitbtc client
	hitbtc := hitbtc.New(API_KEY, API_SECRET)

	// Get balances
	balances, err := hitbtc.GetBalances()
	fmt.Println(err, balances)
}
~~~

In order to use custom settings for the http client do:

~~~ go
package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/sutapurachina/go-hitbtc"
)

const (
	API_KEY    = "YOUR_API_KEY"
	API_SECRET = "YOUR_API_SECRET"
)

func main() {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	// hitbtc client
	bc := hitbtc.NewWithCustomHttpClient(conf.hitbtc.ApiKey, conf.hitbtc.ApiSecret, httpClient)

	// Get balances
	balances, err := hitbtc.GetBalances()
	fmt.Println(err, balances)

	// Initialize websocket connection
	client, err := hitbtc.NewWSClient()
	if err != nil {
		handleError(err) // do something
	}
	defer client.Close()

	// Subscribe and handle
	tickerFeed, err := client.SubscribeTicker("ETHBTC")
	for {
		ticker := <-tickerFeed
		fmt.Println(ticker)
	}


}
~~~

