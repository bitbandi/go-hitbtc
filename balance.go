package hitbtc

// Balance represents a cryptocurrency balance on the exchange
type Balance struct {
	Currency  string  `json:"currency"`
	Available float64 `json:"available,string"`
	Reserved  float64 `json:"reserved,string"`
}
