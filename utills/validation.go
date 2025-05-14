package validation

import (
	"exchange-rate-serivce/enums"
	"fmt"
	"time"
)

func IsValidDate(dateStr string) error {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	// Check if the date is older than 90 days
	if time.Since(date).Hours() > 90*24 {
		return fmt.Errorf("date cannot be older than 90 days")
	}
	return nil
}

func IsValidCurrency(code string) bool {
	return enums.SupportedCurrencies[code]
}
