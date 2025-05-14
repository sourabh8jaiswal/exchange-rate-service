package models

type HistoricateRateResponse struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"amount"`
	Date string  `json:"date"`
}
