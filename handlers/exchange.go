package handlers

import (
	"encoding/json"
	"exchange-rate-serivce/services"
	validation "exchange-rate-serivce/utills"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func Convert(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	if from == "" {
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("from cannot be empty"))
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("to cannot be empty"))
		return
	}

	if !validation.IsValidCurrency(from) || !validation.IsValidCurrency(to) {
		// Check if the currency is valid
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("invalid currency"))
		return
	}

	amount := 1.0
	var err error
	amountStr := r.URL.Query().Get("amount")
	if amountStr != "" {
		amount, err = strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			ErrorWith(w, http.StatusBadRequest, fmt.Errorf("invalid amount"))
			return
		}
	}
	date := r.URL.Query().Get("date")
	if date == "" {
		// if date is not provided, use the current date
		date = time.Now().Format("2006-01-02")
	}

	err = validation.IsValidDate(date)
	if err != nil {
		ErrorWith(w, http.StatusBadRequest, err)
		return
	}

	rate, err := services.GetRate(from, to, date)
	fmt.Println(rate)
	if err != nil {
		ErrorWith(w, http.StatusInternalServerError, err)
		return
	}
	amount = amount * rate

	RespondWith(w, http.StatusOK, &Response{Amount: amount})
}

func LatestRate(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	if from == "" {
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("from cannot be empty"))
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("to cannot be empty"))
		return
	}

	if !validation.IsValidCurrency(from) || !validation.IsValidCurrency(to) {
		// Check if the currency is valid
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("invalid currency"))
		return
	}

	rate, err := services.GetRate(from, to, time.Now().Format("2006-01-02"))
	if err != nil {
		ErrorWith(w, http.StatusInternalServerError, err)
		return
	}
	RespondWith(w, http.StatusOK, &struct{ Rate float64 }{Rate: rate})
}

func HistoricaRates(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	if from == "" {
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("from cannot be empty"))
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("to cannot be empty"))
		return
	}

	if !validation.IsValidCurrency(from) || !validation.IsValidCurrency(to) {
		// Check if the currency is valid
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("invalid currency"))
		return
	}

	sdate := r.URL.Query().Get("startDate")
	if sdate == "" {
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("startDate cannot be empty"))
		return
	}

	err := validation.IsValidDate(sdate)
	if err != nil {
		ErrorWith(w, http.StatusBadRequest, err)
		return
	}

	edate := r.URL.Query().Get("endDate")
	if edate == "" {
		ErrorWith(w, http.StatusBadRequest, fmt.Errorf("endDate cannot be empty"))
		return
	}

	err = validation.IsValidDate(edate)
	if err != nil {
		ErrorWith(w, http.StatusBadRequest, err)
		return
	}

	rates, err := services.GetHistoricalRates(from, to, sdate, edate)
	fmt.Println(rates)
	if err != nil {
		ErrorWith(w, http.StatusInternalServerError, err)
		return
	}

	RespondWith(w, http.StatusOK, &rates)
}

func StartRateRefresher() {
	ticker := time.NewTicker(1 * time.Hour)
	for {
		fmt.Println("Refreshing Exchange Rates!")
		services.FetchAndCacheRates()
		<-ticker.C
	}
}

type Response struct {
	Amount float64 `json:"amount"`
}

func ErrorWith(w http.ResponseWriter, statusCode int, err error) {
	errMsg := struct {
		Error string `json:"message"`
	}{
		Error: err.Error(),
	}

	b, err := json.Marshal(errMsg)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)
}

func RespondWith(w http.ResponseWriter, statusCode int, res interface{}) {
	data, err := json.Marshal(res)
	if err != nil {
		ErrorWith(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
