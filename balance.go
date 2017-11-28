package hitbtc

type Balance struct {
	Currency  string  `json:"currency"`
	Available float64 `json:"available"`
	Reserved  float64 `json:"reserved"`
}
