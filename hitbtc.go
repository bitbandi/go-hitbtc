// Package hitbtc is an implementation of the HitBTC API in Golang.
package hitbtc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	API_BASE = "https://api.hitbtc.com/api/2" // HitBtc API endpoint
)

// New returns an instantiated HitBTC struct
func New(apiKey, apiSecret string) *HitBtc {
	client := NewClient(apiKey, apiSecret)
	return &HitBtc{client}
}

// NewWithCustomHttpClient returns an instantiated HitBTC struct with custom http client
func NewWithCustomHttpClient(apiKey, apiSecret string, httpClient *http.Client) *HitBtc {
	client := NewClientWithCustomHttpConfig(apiKey, apiSecret, httpClient)
	return &HitBtc{client}
}

// NewWithCustomTimeout returns an instantiated HitBTC struct with custom timeout
func NewWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) *HitBtc {
	client := NewClientWithCustomTimeout(apiKey, apiSecret, timeout)
	return &HitBtc{client}
}

// handleErr gets JSON response from livecoin API en deal with error
func handleErr(r interface{}) error {
	switch v := r.(type) {
	case map[string]interface{}:
		error := r.(map[string]interface{})["error"]
		if error != nil {
			switch v := error.(type) {
			case map[string]interface{}:
				errorMessage := error.(map[string]interface{})["message"]
				return errors.New(errorMessage.(string))
			default:
				return fmt.Errorf("I don't know about type %T!\n", v)
			}
		}
	case []interface{}:
		return nil
	default:
		return fmt.Errorf("I don't know about type %T!\n", v)
	}

	return nil
}

// HitBtc represent a HitBTC client
type HitBtc struct {
	client *client
}

// SetDebug sets enable/disable http request/response dump
func (b *HitBtc) SetDebug(enable bool) {
	b.client.debug = enable
}

// GetCurrencies is used to get all supported currencies at HitBtc along with other meta data.
func (b *HitBtc) GetCurrencies() (currencies []Currency, err error) {
	r, err := b.client.do("GET", "public/currency", nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &currencies)
	return
}

// GetSymbols is used to get the open and available trading markets at HitBtc along with other meta data.
func (b *HitBtc) GetSymbols() (symbols []Symbol, err error) {
	r, err := b.client.do("GET", "public/symbol", nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &symbols)
	return
}

// GetTicker is used to get the current ticker values for a market.
func (b *HitBtc) GetTicker(market string) (ticker Ticker, err error) {
	r, err := b.client.do("GET", "public/ticker/"+strings.ToUpper(market), nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &ticker)
	return
}

// GetOrderbook is used to get the current order book for a market.
func (b *HitBtc) GetOrderbook(market string) (orderbook Orderbook, err error) {
	r, err := b.client.do("GET", "public/orderbook/"+strings.ToUpper(market), nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &orderbook)
	return
}

// GetAllTicker is used to get the current ticker values for all markets.
func (b *HitBtc) GetAllTicker() (tickers Tickers, err error) {
	r, err := b.client.do("GET", "public/ticker", nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &tickers)
	return
}

// GetBalances is used to retrieve all balances from your account
func (b *HitBtc) GetBalances() (balances []Balance, err error) {
	r, err := b.client.do("GET", "trading/balance", nil, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &balances)
	return
}

// GetBalance is used to retrieve the balance from your account for a specific currency.
// currency: a string literal for the currency (ex: LTC)
func (b *HitBtc) GetBalance(currency string) (balance Balance, err error) {
	balances, err := b.GetBalances()
	currency = strings.ToUpper(currency)

	for _, balance = range balances {
		if balance.Currency == currency {
			return
		}
	}

	return Balance{}, errors.New("Currency not found")
}

// GetTrades used to retrieve your trade history.
// market string literal for the market (ie. BTC/LTC). If set to "all", will return for all market
func (b *HitBtc) GetTrades(currencyPair string) (trades []Trade, err error) {
	payload := make(map[string]string)
	if currencyPair != "all" {
		payload["symbol"] = currencyPair
	}
	r, err := b.client.do("GET", "history/trades", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &trades)
	return
}

// CancelOrder cancels a pending order
func (b *HitBtc) CancelOrder(currencyPair string) (orders []Order, err error) {
	payload := make(map[string]string)
	if currencyPair != "all" {
		payload["symbol"] = currencyPair
	}
	r, err := b.client.do("DELETE", "order", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &orders)
	return
}

// GetOrder gets a pending order data.
func (b *HitBtc) GetOrder(orderId string) (orders []Order, err error) {
	payload := make(map[string]string)
	payload["clientOrderId"] = orderId
	r, err := b.client.do("GET", "history/order", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &orders)
	return
}

// GetOrderHistory gets the history of orders for an user.
func (b *HitBtc) GetOrderHistory() (orders []Order, err error) {
	r, err := b.client.do("GET", "history/order", nil, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &orders)
	return
}

// GetOpenOrders gets the open orders of an user.
func (b *HitBtc) GetOpenOrders() (orders []Order, err error) {
	r, err := b.client.do("GET", "order", nil, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &orders)
	return
}

// PlaceOrder creates a new order.
func (b *HitBtc) PlaceOrder(requestOrder Order) (responseOrder Order, err error) {
	payload := make(map[string]string)

	payload["symbol"] = requestOrder.Symbol
	payload["side"] = requestOrder.Side
	payload["type"] = requestOrder.Type
	payload["timeInForce"] = requestOrder.TimeInForce
	payload["quantity"] = fmt.Sprintf("%.8f", requestOrder.Quantity)
	payload["price"] = fmt.Sprintf("%.8f", requestOrder.Price)

	r, err := b.client.do("PUT", "order/"+requestOrder.ClientOrderId, payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &responseOrder)
	return
}

// GetTransactions is used to retrieve your withdrawal and deposit history
// "Start" and "end" are given in UNIX timestamp format in miliseconds and used to specify the date range for the data returned.
func (b *HitBtc) GetTransactions(start uint64, end uint64, limit uint32) (transactions []Transaction, err error) {
	payload := make(map[string]string)
	if start > 0 {
		payload["from"] = strconv.FormatUint(uint64(start), 10)
	}
	if end == 0 {
		end = uint64(time.Now().Unix()) * 1000
	}
	if end > 0 {
		payload["till"] = strconv.FormatUint(uint64(end), 10)
	}
	if limit > 1000 {
		limit = 1000
	}
	if limit > 0 {
		payload["limit"] = strconv.FormatUint(uint64(limit), 10)
	}
	r, err := b.client.do("GET", "account/transactions", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	err = json.Unmarshal(r, &transactions)
	return
}

// Withdraw performs a withdrawal operation.
func (b *HitBtc) Withdraw(address string, currency string, amount float64) (withdrawID string, err error) {
	type withdrawResponse struct {
		ID string `json:"id,required"`
	}

	payload := map[string]string{
		"currency": currency,
		"address":  address,
		"amount":   fmt.Sprint(amount),
	}

	r, err := b.client.do("POST", "account/crypto/withdraw", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}

	var withdraw withdrawResponse
	if err = json.Unmarshal(r, &withdraw); err != nil {
		return
	}
	withdrawID = withdraw.ID
	return
}

type transferType string

const (
	// TransferTypeBankToExchange represent a transfer from bank (withdraw) balance to exchange (trading) balance.
	TransferTypeBankToExchange transferType = "bankToExchange"
	// TransferTypeExchangeToBank represent a transfer from exchange (trading) balance to bank (withdraw) balance.
	TransferTypeExchangeToBank transferType = "exchangeToBank"
)

// TransferBalance performs a balance transfer operation between trading and bank accounts (both directions).
func (b *HitBtc) TransferBalance(currency string, amount float64, transferType transferType) (transferID string, err error) {
	type transferResponse struct {
		ID string `json:"id,required"`
	}

	payload := map[string]string{
		"currency": currency,
		"amount":   fmt.Sprint(amount),
		"type":     string(transferType),
	}

	r, err := b.client.do("POST", "account/transfer", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}

	var transfer transferResponse
	if err = json.Unmarshal(r, &transfer); err != nil {
		return
	}
	transferID = transfer.ID
	return
}
