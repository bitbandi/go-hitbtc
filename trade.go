package hitbtc

import (
	"encoding/json"
	"time"
)

type Trade struct {
	Id            uint64    `json:"id"`
	OrderId       uint64    `json:"orderId"`
	ClientOrderId string    `json:"clientOrderId"`
	Symbol        string    `json:"symbol"`
	Type          string    `json:"side"`
	Price         float64   `json:"price,string"`
	Quantity      float64   `json:"quantity,string"`
	Fee           float64   `json:"fee,string"`
	Timestamp     time.Time `json:"timestamp"`
}

func (t *Trade) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Trade
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
