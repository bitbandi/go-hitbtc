package hitbtc

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/juju/errors"
	jsonrpc2 "github.com/sourcegraph/jsonrpc2"
	jsonrpc2ws "github.com/sourcegraph/jsonrpc2/websocket"
)

const wsAPIURL string = "wss://api.hitbtc.com/api/2/ws"

type noopHandler struct{}

func (noopHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {}

var noopHandle = noopHandler{}

// WSClient represents a JSON RPC v2 Connection over Websocket,
type WSClient struct {
	conn *jsonrpc2.Conn
}

// NewWSClient creates a new WSClient
func NewWSClient() (*WSClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsAPIURL, nil)
	if err != nil {
		return nil, err
	}

	return &WSClient{
		conn: jsonrpc2.NewConn(context.Background(), jsonrpc2ws.NewObjectStream(conn), noopHandle),
	}, nil
}

// Close closes the Websocket connected to the hitbtc api.
func (c *WSClient) Close() {
	c.conn.Close()
}

// WSGetCurrencyRequest is get currency request type on websocket
type WSGetCurrencyRequest struct {
	Currency string `json:"currency,required"`
}

// WSGetCurrencyResponse is get currency response type on websocket
type WSGetCurrencyResponse struct {
	ID                 string `json:"id,required"`
	FullName           string `json:"fullname,required"`
	Crypto             bool   `json:"crypto,required"`
	PayinEnabled       bool   `json:"payinEnabled,required"`
	PayinPaymentID     bool   `json:"payinPaymentId,required"`
	PayinConfirmations int    `json:"payinConfirmations,required"`
	PayoutEnabled      bool   `json:"payoutEnabled,required"`
	PayoutIsPaymentID  bool   `json:"payoutIsPaymentId,required"`
	TransferEnabled    bool   `json:"transferEnabled,required"`
	Delisted           bool   `json:"delisted,required"`
	PayoutFee          string `json:"payoutFee,required"`
}

// GetCurrencyInfo get the info about a currency.
func (c *WSClient) GetCurrencyInfo(symbol string) (*WSGetCurrencyResponse, error) {
	var request = WSGetCurrencyRequest{Currency: symbol}
	var response WSGetCurrencyResponse

	err := c.conn.Call(context.Background(), "getCurrency", request, &response)
	if err != nil {
		return nil, errors.Annotate(err, "Hitbtc GetCurrency")
	}
	return &response, nil
}

// WSGetSymbolsRequest is get symbols request type on websocket
type WSGetSymbolsRequest struct {
	Symbol string `json:"symbol,required"`
}

// WSGetSymbolsResponse  is get symbols response type on websocket
type WSGetSymbolsResponse struct {
	ID                   string `json:"id,required"`
	BaseCurrency         string `json:"baseCurrency,required"`
	QuoteCurrency        string `json:"quoteCurrency,required"`
	QuantityIncrement    string `json:"quantityIncrement,required"`
	TickSize             string `json:"tickSize,required"`
	TakeLiquidityRate    string `json:"takeLiquidityRate,required"`
	ProvideLiquidityRate string `json:"provideLiquidityRate,required"`
	FeeCurrency          string `json:"feeCurrency,required"`
}

// WSSubscribeTickerRequest is get symbols request type on websocket
type WSSubscribeTickerRequest struct {
	Symbol string `json:"symbol,required"`
}

// WSNotificationTickerResponse is notification response type on websocket
type WSNotificationTickerResponse struct {
	Ask         string `json:"ask,required"`         // Best ask price
	Bid         string `json:"bid,required"`         // Best bid price
	Last        string `json:"last,required"`        // Last trade price
	Open        string `json:"open,required"`        // Last trade price 24 hours ago
	Low         string `json:"low,required"`         // Lowest trade price within 24 hours
	High        string `json:"high,required"`        // Highest trade price within 24 hours
	Volume      string `json:"volume,required"`      // Total trading amount within 24 hours in base currency
	VolumeQuote string `json:"volumeQuote,required"` // Total trading amount within 24 hours in quote currency
	Timestamp   string `json:"timestamp,required"`   // Last update or refresh ticker timestamp
	Symbol      string `json:"symbol,required"`
}

// WSSubscribeOrderbookRequest is subscribe request type to orderbook on websocket
type WSSubscribeOrderbookRequest struct {
	Symbol string `json:"symbol,required"`
}

// WSSubtypeTrade is element of market trade type
type WSSubtypeTrade struct {
	Price string `json:"price,required"`
	Size  string `json:"size,required"`
}

// WSNotificationOrderbookSnapshot is notification response type to orderbook snapshot on websocket
type WSNotificationOrderbookSnapshot struct {
	Ask      []WSSubtypeTrade `json:"ask,required"`
	Bid      []WSSubtypeTrade `json:"bid,required"`
	Symbol   string           `json:"symbol,required"`
	Sequence int64            `json:"sequence,required"`
}

// WSNotificationOrderbookUpdate is notification response type to orderbook snapshot on websocket
type WSNotificationOrderbookUpdate struct {
	Ask      []WSSubtypeTrade `json:"ask,required"`
	Bid      []WSSubtypeTrade `json:"bid,required"`
	Symbol   string           `json:"symbol,required"`
	Sequence int64            `json:"sequence,required"`
}

// WSSubscribeTradesRequest is subscribe request type to trades on websocket
type WSSubscribeTradesRequest struct {
	Symbol string `json:"symbol,required"`
}

// WSTrades is item for Trades
type WSTrades struct {
	ID        int    `json:"id,required"`
	Price     string `json:"price,required"`
	Quantity  string `json:"quantity"`
	Side      string `json:"side,required"`
	Timestamp string `json:"timestamp,required"`
}

// WSNotificationTradesSnapshot is notification response type to trades on websocket
type WSNotificationTradesSnapshot struct {
	Data []WSTrades `json:"data,required"`
}

// WSNotificationTradesUpdate is notification response type to trades on websocket
type WSNotificationTradesUpdate struct {
	Data WSTrades `json:"data,required"`
}

// WSGetTradesRequest is get trades request type on websocket
type WSGetTradesRequest struct {
	Symbol string    `json:"symbol,required"`
	Limit  int       `json:"limit,required"`
	Sort   string    `json:"sort,required"`
	By     string    `json:"by,required"`
	From   time.Time `json:"from"`
	Till   time.Time `json:"till"`
	Offset string    `json:"offset"`
}

// WSGetTradesResponse  is get symbols response type on websocket
type WSGetTradesResponse struct {
	Data []WSTrades `json:"data,required"`
}

// WSSubscribeCandlesRequest is subscribe request type to candles on websocket
type WSSubscribeCandlesRequest struct {
	Symbol string `json:"symbol,required"`
	Period string `json:"period,required"`
}

// WSCandles is item for WSCandles
type WSCandles struct {
	Timestamp   time.Time `json:"timestamp,required"`
	Open        string    `json:"open,required"`
	Close       string    `json:"close,required"`
	Min         string    `json:"min,required"`
	Max         string    `json:"max,required"`
	Volume      string    `json:"volume,required"`      // Total trading amount within 24 hours in base currency
	VolumeQuote string    `json:"volumeQuote,required"` // Total trading amount within 24 hours in quote currency
}

// WSSubscribeCandlesNotificationSnapshot is subscribe response type to candles on websocket
type WSSubscribeCandlesNotificationSnapshot struct {
	Data   []WSCandles `json:"data,required"`
	Symbol string      `json:"symbol,required"`
	Period string      `json:"period,required"`
}

// WSSubscribeCandlesNotificationUpdate is subscribe response type to candles on websocket
type WSSubscribeCandlesNotificationUpdate struct {
	Data   WSCandles `json:"data,required"`
	Symbol string    `json:"symbol,required"`
	Period string    `json:"period,required"`
}
