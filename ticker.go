package hitbtc

import (
	"encoding/json"
	"time"
)

//Ticker represents a Ticker from hitbtc API.
type Ticker struct {
	Ask         float64   `json:"ask,string"`
	Bid         float64   `json:"bid,string"`
	Last        float64   `json:"last,string"`
	Open        float64   `json:"open,string"`
	Low         float64   `json:"low,string"`
	High        float64   `json:"high,string"`
	Volume      float64   `json:"volume,string"`
	VolumeQuote float64   `json:"volumeQuote,string"`
	Timestamp   time.Time `json:"timestamp"`
	Symbol      string    `json:"symbol"`
}

// Tickers rapresents a set of a valid Tickers struct
type Tickers []Ticker

func (t *Ticker) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Ticker
	aux := &struct {
		Timestamp string `json:"timestamp"`
		High      string `json:"high"`
		Low       string `json:"low"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if aux.High == "" {
		aux.High = "0"
	}
	if aux.Low == "" {
		aux.Low = "0"
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

func (ts Tickers) UnmarshalJSON(data []byte) error {
	var arr []json.RawMessage

	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	for _, v := range arr {
		var ticker Ticker

		err := json.Unmarshal(v, &ticker)
		if ticker.Volume == 0 {
			continue
		}
		if err != nil {
			return err
		}
		ts = append(ts, ticker)
	}

	return nil
}
