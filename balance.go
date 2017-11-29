package hitbtc

type Balance struct {
	Currency  string  `json:"currency"`
	Available float64 `json:"available,string"`
	Reserved  float64 `json:"reserved,string"`
}
