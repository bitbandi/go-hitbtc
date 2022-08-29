package hitbtc_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/sutapurachina/go-hitbtc"
)

const (
	// email: g3217755@nwytg.net
	// password: Test10
	apiKey    = "7567417ba8df166f50584b54ccd924c5"
	apiSecret = "2203d97807b478326c4bcade686538b0"
)

var (
	hitBtc              = hitbtc.New(apiKey, apiSecret)
	defaultErrorMessage = "There should be no error"
)

func TestGetCurrencies(t *testing.T) {
	currencies, err := hitBtc.GetCurrencies()
	t.Logf("GetCurrencies : %#v\n", currencies)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetSymbols(t *testing.T) {
	symbols, err := hitBtc.GetSymbols()
	t.Logf("GetSymbols : %#v\n", symbols)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetTicker(t *testing.T) {
	ticker, err := hitBtc.GetTicker("ETHBTC")
	t.Logf("GetTicker : %#v\n", ticker)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetOrderbook(t *testing.T) {
	orderbook, err := hitBtc.GetOrderbook("ETHBTC")
	t.Logf("GetOrderbook : %#v\n", orderbook)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetAllTicker(t *testing.T) {
	tickers, err := hitBtc.GetAllTicker()
	t.Logf("GetAllTicker : %v\n", tickers)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetBalances(t *testing.T) {
	balances, err := hitBtc.GetBalances()
	t.Logf("GetBalances : %#v\n", balances)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetTrades(t *testing.T) {
	trades, err := hitBtc.GetTrades("ETHBTC")
	t.Logf("GetTrades : %#v\n", trades)
	require.NoError(t, err, defaultErrorMessage)
}

func TestCancelOrder(t *testing.T) {
	orders, err := hitBtc.CancelOrder("ETHBTC")
	t.Logf("CancelOrder : %#v\n", orders)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetOrder(t *testing.T) {
	orders, err := hitBtc.GetOrder("ETHBTC")
	t.Logf("GetOrder : %#v\n", orders)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetOrderHistory(t *testing.T) {
	orders, err := hitBtc.GetOrderHistory()
	t.Logf("GetOrderHistory : %#v\n", orders)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetOpenOrders(t *testing.T) {
	orders, err := hitBtc.GetOpenOrders()
	t.Logf("GetOpenOrders : %#v\n", orders)
	require.NoError(t, err, defaultErrorMessage)
}
