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

// responseChannels handles all incoming data from the hitbtc connection.
type responseChannels struct {
	notifications notificationChannels

	OrderbookFeed chan WSNotificationOrderbookSnapshot
	TradesFeed    chan WSNotificationTradesSnapshot
	CandlesFeed   chan WSNotificationCandlesSnapshot

	ErrorFeed chan error
}

// notificationChannels contains all the notifications from hitbtc for subscribed feeds.
type notificationChannels struct {
	TickerFeed    chan WSNotificationTickerResponse
	OrderbookFeed chan WSNotificationOrderbookUpdate
	TradesFeed    chan WSNotificationTradesUpdate
	CandlesFeed   chan WSNotificationCandlesUpdate
}

// Handle handles all incoming connections and fills the channels properly.
func (h *responseChannels) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	if req.Params != nil {
		message := *req.Params
		switch req.Method {
		case "ticker":
			var msg WSNotificationTickerResponse
			err := json.Unmarshal(message, &msg)
			if err != nil {
				h.ErrorFeed <- err
			} else {
				h.notifications.TickerFeed <- msg
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
				h.notifications.OrderbookFeed <- msg
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
				h.notifications.TradesFeed <- msg
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
				h.notifications.CandlesFeed <- msg
			}
		}
	}
}

// WSClient represents a JSON RPC v2 Connection over Websocket,
type WSClient struct {
	conn    *jsonrpc2.Conn
	updates *responseChannels
}

// NewWSClient creates a new WSClient
func NewWSClient() (*WSClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsAPIURL, nil)
	if err != nil {
		return nil, err
	}

	var handler responseChannels

	return &WSClient{
		conn:    jsonrpc2.NewConn(context.Background(), jsonrpc2ws.NewObjectStream(conn), &handler),
		updates: &handler,
	}, nil
}

// Close closes the Websocket connected to the hitbtc api.
func (c *WSClient) Close() {
	c.conn.Close()
	close(c.updates.notifications.TickerFeed)
	close(c.updates.notifications.OrderbookFeed)
	close(c.updates.notifications.TradesFeed)
	close(c.updates.notifications.CandlesFeed)
	close(c.updates.OrderbookFeed)
	close(c.updates.TradesFeed)
	close(c.updates.CandlesFeed)
	close(c.updates.ErrorFeed)
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

// WSGetTradesRequest is get trades request type on websocket
type WSGetTradesRequest struct {
	Symbol string     `json:"symbol,required"`
	Limit  int        `json:"limit,required"`
	Sort   string     `json:"sort,required"`
	By     string     `json:"by,required"`
	From   *time.Time `json:"from,omitempty"`
	Till   *time.Time `json:"till,omitempty"`
	Offset *string    `json:"offset,omitempty"`
}

// WSGetTradesResponse  is get symbols response type on websocket
type WSGetTradesResponse struct {
	Data []WSTrades `json:"data,required"`
}

// GetTrades obtains the data of a series of trades, based on the specified filters.
func (c *WSClient) GetTrades(symbol string) (*WSGetTradesResponse, error) {
	var request = WSGetTradesRequest{Symbol: symbol}
	var response WSGetTradesResponse

	err := c.conn.Call(context.Background(), "getSymbol", request, &response)
	if err != nil {
		return nil, errors.Annotate(err, "Hitbtc GetSymbol")
	}
	return &response, nil
}

// wsSubscriptionResponse is the response for a subscribe/unsubscribe requests.
type wsSubscriptionResponse struct {
	Result bool `json:"result,required"`
}

// WSSubscriptionRequest is request type on websocket subscription.
type WSSubscriptionRequest struct {
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

// SubscribeTicker subscribes to the specified market ticker notifications.
func (c *WSClient) SubscribeTicker(symbol string) (<-chan WSNotificationTickerResponse, error) {
	err := c.subscriptionOp("subscribeTicker", symbol)
	if err != nil {
		return nil, errors.Annotate(err, "Hitbtc SubscribeTicker")
	}

	c.updates.notifications.TickerFeed = make(chan WSNotificationTickerResponse)

	return c.updates.notifications.TickerFeed, nil
}

// UnsubscribeTicker subscribes to the specified market ticker notifications.
//
// This closes also the connected channel of updates.
func (c *WSClient) UnsubscribeTicker(symbol string) error {
	err := c.subscriptionOp("unsubscribeTicker", symbol)
	if err != nil {
		return errors.Annotate(err, "Hitbtc UnsubscribeTicker")
	}

	close(c.updates.notifications.TickerFeed)

	return nil
}

// WSNotificationTradesSnapshot is notification response type to trades on websocket
type WSNotificationTradesSnapshot struct {
	Data []WSTrades `json:"data,required"`
}

// WSNotificationTradesUpdate is notification response type to trades on websocket
type WSNotificationTradesUpdate struct {
	Data WSTrades `json:"data,required"`
}

// WSTrades is item for Trades
type WSTrades struct {
	ID        int    `json:"id,required"`
	Price     string `json:"price,required"`
	Quantity  string `json:"quantity"`
	Side      string `json:"side,required"`
	Timestamp string `json:"timestamp,required"`
}

// SubscribeTrades subscribes to the specified market trades notifications.
func (c *WSClient) SubscribeTrades(symbol string) (<-chan WSNotificationTradesUpdate, <-chan WSNotificationTradesSnapshot, error) {
	err := c.subscriptionOp("subscribeTrades", symbol)
	if err != nil {
		return nil, nil, errors.Annotate(err, "Hitbtc SubscribeTrades")
	}

	c.updates.notifications.TradesFeed = make(chan WSNotificationTradesUpdate)
	c.updates.TradesFeed = make(chan WSNotificationTradesSnapshot)

	return c.updates.notifications.TradesFeed, c.updates.TradesFeed, nil
}

// UnsubscribeTrades unsubscribes from the specified market trades notifications and snapshot.
//
// This closes also the connected channel of updates.
func (c *WSClient) UnsubscribeTrades(symbol string) error {
	err := c.subscriptionOp("unsubscribeTrades", symbol)
	if err != nil {
		return errors.Annotate(err, "Hitbtc UnsubscribeTrades")
	}

	close(c.updates.notifications.TradesFeed)
	close(c.updates.TradesFeed)

	return nil
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

// SubscribeOrderbook subscribes to the specified market order book notifications.
func (c *WSClient) SubscribeOrderbook(symbol string) (<-chan WSNotificationOrderbookUpdate, <-chan WSNotificationOrderbookSnapshot, error) {
	err := c.subscriptionOp("subscribeOrderbook", symbol)
	if err != nil {
		return nil, nil, errors.Annotate(err, "Hitbtc SubscribeOrderbook")
	}

	c.updates.notifications.OrderbookFeed = make(chan WSNotificationOrderbookUpdate)
	c.updates.OrderbookFeed = make(chan WSNotificationOrderbookSnapshot)

	return c.updates.notifications.OrderbookFeed, c.updates.OrderbookFeed, nil
}

// UnsubscribeOrderbook unsubscribes from the specified market order book notifications and snapshot.
//
// This closes also the connected channel of updates.
func (c *WSClient) UnsubscribeOrderbook(symbol string) error {
	err := c.subscriptionOp("unsubscribeOrderbook", symbol)
	if err != nil {
		return errors.Annotate(err, "Hitbtc UnsubscribeOrderbook")
	}

	close(c.updates.notifications.OrderbookFeed)
	close(c.updates.OrderbookFeed)

	return nil
}

const (
	// Interval30Minutes is 30 minutes interval for candle data.
	Interval30Minutes string = "M30"
	// Interval1Hour is 1 hour interval for candle data.
	Interval1Hour string = "H1"
)

// WSCandlesSubscriptionRequest is a request to subscribe for candle data.
type WSCandlesSubscriptionRequest struct {
	Symbol string `json:"symbol,required"`
	Period string `json:"period,required"`
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

// SubscribeCandles subscribes to the specified market candle notifications for the specified timeframe.
func (c *WSClient) SubscribeCandles(symbol string, timeframe string) (<-chan WSNotificationCandlesUpdate, <-chan WSNotificationCandlesSnapshot, error) {
	err := c.candlesSubscriptionOp("subscribeCandles", symbol, timeframe)
	if err != nil {
		return nil, nil, errors.Annotate(err, "Hitbtc SubscribeCandles")
	}

	c.updates.notifications.CandlesFeed = make(chan WSNotificationCandlesUpdate)
	c.updates.CandlesFeed = make(chan WSNotificationCandlesSnapshot)

	return c.updates.notifications.CandlesFeed, c.updates.CandlesFeed, nil
}

// UnsubscribeCandles unsubscribes from the specified market candle notifications for the specified timeframe.
//
// This closes also the connected channel of updates.
func (c *WSClient) UnsubscribeCandles(symbol string, timeframe string) error {
	err := c.candlesSubscriptionOp("unsubscribeCandles", symbol, timeframe)
	if err != nil {
		return errors.Annotate(err, "Hitbtc UnsubscribeCandles")
	}

	close(c.updates.notifications.CandlesFeed)
	close(c.updates.CandlesFeed)

	return nil
}

func (c *WSClient) subscriptionOp(op string, symbol string) error {
	var request = WSSubscriptionRequest{Symbol: symbol}
	var response wsSubscriptionResponse

	err := c.conn.Call(context.Background(), op, request, &response)
	if err != nil {
		return err
	}

	return nil
}

func (c *WSClient) candlesSubscriptionOp(op string, symbol string, period string) error {
	var request = WSCandlesSubscriptionRequest{Symbol: symbol, Period: period}
	var response wsSubscriptionResponse

	err := c.conn.Call(context.Background(), op, request, &response)
	if err != nil {
		return err
	}

	return nil
}
