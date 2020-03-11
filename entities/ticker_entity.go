package entities

type Ticker struct {
	ID     int `json:"id"`
	Symbol string `json:"symbol"`
	Value  float64 `json:"value"`
	Quota float64 `json:"quota"`
	AvgPrice float64 `json:"avgPrice"`
	PreviousClose float64 `json:"previousClose"`
	LastChangePercent float64 `json:"lastChangePercent"`
	ChangeFromAvgPrice float64 `json:"changeFromAvgPrice"`
}
