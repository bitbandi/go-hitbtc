package hitbtc

import (
    "encoding/json"
    "time"
)

type Order struct {
    Id            string    `json:"id"`
    ClientOrderId string    `json:"clientOrderId"`
    Symbol        string    `json:"symbol"`
    Side          string    `json:"side"`
    Status        string    `json:"status"`
    Type          string    `json:"type"`
    TimeInForce   string    `json:"timeInForce"`
    Quantity      float64   `json:"quantity,string"`
    Price         float64   `json:"price,string"`
    CumQuantity   float64   `json:"cumQuantity,string"`
    Created       time.Time `json:"createdAt"`
    Updated       time.Time `json:"updatedAt"`
    StopPrice     float64   `json:"stopPrice,string"`
    Expire        time.Time `json:"expireTime"`
}

func (t *Order) UnmarshalJSON(data []byte) error {
    var err error
    type Alias Order
    aux := &struct {
        Created string `json:"createdAt"`
        Updated string `json:"updatedAt"`
        Expire  string `json:"expireTime"`
        *Alias
    }{
        Alias: (*Alias)(t),
    }
    if err = json.Unmarshal(data, &aux); err != nil {
        return err
    }
    if aux.Created != "" {
        t.Created, err = time.Parse("2006-01-02T15:04:05.999Z", aux.Created)
        if err != nil {
            return err
        }
    }
    if aux.Updated != "" {
        t.Updated, err = time.Parse("2006-01-02T15:04:05.999Z", aux.Updated)
        if err != nil {
            return err
        }
    }
    if aux.Expire != "" {
        t.Expire, err = time.Parse("2006-01-02T15:04:05.999Z", aux.Expire)
        if err != nil {
            return err
        }
    }
    return nil
}
