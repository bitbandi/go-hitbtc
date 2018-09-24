package hitbtc

import (
	"encoding/json"
	"time"
)

// Orderbook represents an orderbook from hitbtc api.
type Orderbook struct {
	Ask       []OrederBookItem `json:"ask,struct"`
	Bid       []OrederBookItem `json:"bid,struct"`
	Timestamp time.Time        `json:"timestamp"`
}

// OrederBookItem for Ask and Bid field
type OrederBookItem struct {
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
	t.Timestamp, err = time.Parse("2006-01-02T15:04:05.999Z", aux.Timestamp)
	if err != nil {
		return err
	}
	return nil
}