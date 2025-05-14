package services

import (
	"encoding/json"
	"exchange-rate-serivce/cache"
	"exchange-rate-serivce/enums"
	"exchange-rate-serivce/models"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const API_URL = "http://api.exchangerate.host/convert"
const ACCESS_KEY = "277229723a34baced2970179b15f70b8"

func GetHistoricalRates(from string, to string, sdate string, edate string) ([]models.HistoricateRateResponse, error) {
	var result []models.HistoricateRateResponse
	layout := "2006-01-02" // Go's reference layout for time formatting
	startDate, err := time.Parse(layout, sdate)
	if err != nil {
		return nil, fmt.Errorf("error in parsing startDate: %v", err)
	}
	endDate, err := time.Parse(layout, edate)
	if err != nil {
		return nil, fmt.Errorf("error in parsing endDate: %v", err)

	}

	// Loop over each day between startDate and endDate
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		strDate := d.Format(layout)
		rate, err := GetRate(from, to, strDate)
		if err != nil {
			return result, err
		}
		result = append(result, models.HistoricateRateResponse{
			From: from,
			To:   to,
			Date: strDate,
			Rate: rate,
		})

	}

	return result, nil
}

func GetRate(from string, to string, date string) (float64, error) {

	r := cache.GetRate(from, to, date)
	if r != 0 {
		return r, nil
	}

	query := url.Values{}
	query.Add("access_key", ACCESS_KEY)
	query.Add("from", from)
	query.Add("to", to)
	query.Add("amount", "1")
	query.Add("date", date)

	url := API_URL + "?" + query.Encode()

	httpRes, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("network err: %v", err)
	}

	bodyBytes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return 0, err
	}

	if httpRes.StatusCode != http.StatusOK {
		err = fmt.Errorf("api error: %s", string(bodyBytes))
		return 0, err
	}

	fmt.Println(string(bodyBytes))

	var res CurrencyResponse
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return 0, err
	}
	cache.StoreRate(from, to, date, res.Result)
	return res.Result, nil
}

func FetchAndCacheRates() {
	today := time.Now().Format("2006-01-02")
	for from := range enums.SupportedCurrencies {
		for to := range enums.SupportedCurrencies {
			if from == to {
				continue
			}
			time.Sleep(1 * time.Second)
			GetRate(from, to, today) // Force fetch + cache
		}
	}
}

type CurrencyResponse struct {
	Success    bool    `json:"success"`
	Terms      string  `json:"terms"`
	Privacy    string  `json:"privacy"`
	Query      Query   `json:"query"`
	Info       Info    `json:"info"`
	Historical bool    `json:"historical"`
	Date       string  `json:"date"`
	Result     float64 `json:"result"`
}

type Query struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type Info struct {
	Timestamp int64   `json:"timestamp"`
	Quote     float64 `json:"quote"`
}
