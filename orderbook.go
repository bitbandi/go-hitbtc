package hitbtc

import (
	"encoding/json"
)

// Orderbook represents an orderbook from hitbtc api.
type Orderbook struct {
	Ask []OrderBookItem `json:"ask,struct"`
	Bid []OrderBookItem `json:"bid,struct"`
}

// OrderBookItem for Ask and Bid field.
type OrderBookItem struct {
	Price float64 `json:"price,string"`
	Size  float64 `json:"size,string"`
}

// UnmarshalJSON for OrderBook function
func (t *Orderbook) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Orderbook
	aux := &struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
