package entities

type Ticker struct {
	ID     int `json:"id"`
	Symbol string `json:"symbol"`
	Value  float64 `json:"value"`
	Quota float64 `json:"quota"`
}
