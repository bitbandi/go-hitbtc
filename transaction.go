package hitbtc

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	Id         string    `json:"id"`
	Index      uint64    `json:"index"`
	Currency   string    `json:"currency"`
	Amount     float64   `json:"amount,string"`
	Fee        float64   `json:"fee,string"`
	NetworkFee float64   `json:"networkFee,string"`
	Address    string    `json:"address"`
	Hash       string    `json:"hash"`
	Status     string    `json:"status"`
	Type       string    `json:"type"`
	Created    time.Time `json:"createdAt"`
	Updated    time.Time `json:"updatedAt"`
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Transaction
	aux := &struct {
		Created string `json:"createdAt"`
		Updated string `json:"updatedAt"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Created, err = time.Parse("2006-01-02T15:04:05.999Z", aux.Created)
	if err != nil {
		return err
	}
	t.Updated, err = time.Parse("2006-01-02T15:04:05.999Z", aux.Updated)
	if err != nil {
		return err
	}
	return nil
}
