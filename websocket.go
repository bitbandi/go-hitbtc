package hitbtc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/juju/errors"
	jsonrpc2 "github.com/sourcegraph/jsonrpc2"
	jsonrpc2ws "github.com/sourcegraph/jsonrpc2/websocket"
)

const wsAPIURL string = "wss://api.hitbtc.com/api/2/ws"

// ResponseChannels handles all incoming data from the hitbtc connection.
type ResponseChannels struct {
	Notifications NotificationChannels

	OrderbookFeed chan WSNotificationOrderbookSnapshot
	TradesFeed    chan WSNotificationTradesSnapshot
	CandlesFeed   chan WSNotificationCandlesSnapshot

	ErrorFeed chan error
}

// NotificationChannels contains all the notifications from hitbtc for subscribed feeds.
type NotificationChannels struct {
	TickerFeed    chan WSNotificationTickerResponse
	OrderBookFeed chan WSNotificationOrderbookUpdate
	TradesFeed    chan WSNotificationTradesUpdate
	CandlesFeed   chan WSNotificationCandlesUpdate
}

// Handle handles all incoming connections and fills the channels properly.
func (h *ResponseChannels) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	if req.Params != nil {
		message := *req.Params
		switch req.Method {
		case "ticker":
			var msg WSNotificationTickerResponse
			err := json.Unmarshal(message, &msg)
			if err != nil {
				h.ErrorFeed <- err
			} else {
				h.Notifications.TickerFeed <- msg
			}
		case "snapshotOrderbook":
			var msg WSNotificationOrderbookSnapshot
			err := json.Unmarshal(message, &msg)
			if err != nil {
				h.ErrorFeed <- err
			} else {
				h.OrderbookFeed <- msg
			}
		case "updateOrderbook":
			var msg WSNotificationOrderbookUpdate
			err := json.Unmarshal(message, &msg)
			if err != nil {
				h.ErrorFeed <- err
			} else {
				h.Notifications.OrderBookFeed <- msg
			}
		case "snapshotTrades":
			var msg WSNotificationTradesSnapshot
			err := json.Unmarshal(message, &msg)
			if err != nil {
				h.ErrorFeed <- err
			} else {
				h.TradesFeed <- msg
			}
		case "updateTrades":
			var msg WSNotificationTradesUpdate
			err := json.Unmarshal(message, &msg)
			if err != nil {
				h.ErrorFeed <- err
			} else {
				h.Notifications.TradesFeed <- msg
			}
		case "snapshotCandles":
			var msg WSNotificationCandlesSnapshot
			err := json.Unmarshal(message, &msg)
			if err != nil {
				h.ErrorFeed <- err
			} else {
				h.CandlesFeed <- msg
			}
		case "updateCandles":
			var msg WSNotificationCandlesUpdate
			err := json.Unmarshal(message, &msg)
			if err != nil {
				h.ErrorFeed <- err
			} else {
				h.Notifications.CandlesFeed <- msg
			}
		}
	}
}

// WSClient represents a JSON RPC v2 Connection over Websocket,
type WSClient struct {
	conn    *jsonrpc2.Conn
	Updates *ResponseChannels
}

// NewWSClient creates a new WSClient
func NewWSClient() (*WSClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsAPIURL, nil)
	if err != nil {
		return nil, err
	}

	handler := &ResponseChannels{
		Notifications: NotificationChannels{
			TickerFeed:    make(chan WSNotificationTickerResponse),
			OrderBookFeed: make(chan WSNotificationOrderbookUpdate),
			TradesFeed:    make(chan WSNotificationTradesUpdate),
			CandlesFeed:   make(chan WSNotificationCandlesUpdate),
		},
		OrderbookFeed: make(chan WSNotificationOrderbookSnapshot),
		TradesFeed:    make(chan WSNotificationTradesSnapshot),
		CandlesFeed:   make(chan WSNotificationCandlesSnapshot),
		ErrorFeed:     make(chan error),
	}

	return &WSClient{
		conn:    jsonrpc2.NewConn(context.Background(), jsonrpc2ws.NewObjectStream(conn), handler),
		Updates: handler,
	}, nil
}

// Close closes the Websocket connected to the hitbtc api.
func (c *WSClient) Close() {
	c.conn.Close()
	close(c.Updates.Notifications.TickerFeed)
	close(c.Updates.Notifications.OrderBookFeed)
	close(c.Updates.Notifications.TradesFeed)
	close(c.Updates.Notifications.CandlesFeed)
	close(c.Updates.OrderbookFeed)
	close(c.Updates.TradesFeed)
	close(c.Updates.CandlesFeed)
	close(c.Updates.ErrorFeed)
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

// WSGetSymbolRequest is get symbols request type on websocket
type WSGetSymbolRequest struct {
	Symbol string `json:"symbol,required"`
}

// WSGetSymbolResponse is get symbols response type on websocket
type WSGetSymbolResponse struct {
	ID                   string `json:"id,required"`
	BaseCurrency         string `json:"baseCurrency,required"`
	QuoteCurrency        string `json:"quoteCurrency,required"`
	QuantityIncrement    string `json:"quantityIncrement,required"`
	TickSize             string `json:"tickSize,required"`
	TakeLiquidityRate    string `json:"takeLiquidityRate,required"`
	ProvideLiquidityRate string `json:"provideLiquidityRate,required"`
	FeeCurrency          string `json:"feeCurrency,required"`
}

// GetSymbol obtains the data of a market.
func (c *WSClient) GetSymbol(symbol string) (*WSGetSymbolResponse, error) {
	var request = WSGetSymbolRequest{Symbol: symbol}
	var response WSGetSymbolResponse

	err := c.conn.Call(context.Background(), "getSymbol", request, &response)
	if err != nil {
		return nil, errors.Annotate(err, "Hitbtc GetSymbol")
	}
	return &response, nil
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
	Sequence int64            `json:"sequence,required"` // used to see if update is the latest received
}

// WSNotificationOrderbookUpdate is notification response type to orderbook snapshot on websocket
type WSNotificationOrderbookUpdate struct {
	Ask      []WSSubtypeTrade `json:"ask,required"`
	Bid      []WSSubtypeTrade `json:"bid,required"`
	Symbol   string           `json:"symbol,required"`
	Sequence int64            `json:"sequence,required"` // used to see if the snapshot is the latest
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

// WSNotificationCandlesSnapshot is subscribe response type to candles on websocket
type WSNotificationCandlesSnapshot struct {
	Data   []WSCandles `json:"data,required"`
	Symbol string      `json:"symbol,required"`
	Period string      `json:"period,required"`
}

// WSNotificationCandlesUpdate is subscribe response type to candles on websocket
type WSNotificationCandlesUpdate struct {
	Data   WSCandles `json:"data,required"`
	Symbol string    `json:"symbol,required"`
	Period string    `json:"period,required"`
}
