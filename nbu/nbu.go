package nbu

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Currency struct {
	StartDate     string
	TimeSign      string
	CurrencyCode  string
	CurrencyCodeL string
	Units         int
	Amount        float64
}

// HTTPClient interface
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

var Client HTTPClient

// added to make mock in tests
func init() {
	Client = &http.Client{}
}

// define local rate map to avoid excessive NBU API call for the same date rate
var dateRateMap sync.Map //map[string]float64{}

// Get currency rate to UAH on specified date
func GetConversionRates(date time.Time, currency string) (float64, error) {
	nbuDateFormatString := date.Format("02012006")
	// check if Rate is already present in local rate map
	if v, ok := dateRateMap.Load(fmt.Sprintf("%s/%s", currency, nbuDateFormatString)); ok {
		//fmt.Printf("Currency rate was found in local map %f\n",val)
		val := v.(float64)
		return val, nil
	} else {
		// If rate not in local map make request to NBU rate API
		currencyRates := make([]Currency, 0)
		url := "https://bank.gov.ua/NBU_Exchange/exchange?json&date=" + nbuDateFormatString

		resp, err := Client.Get(url)
		if err != nil {
			return float64(0), err
		}
		defer resp.Body.Close()

		json.NewDecoder(resp.Body).Decode(&currencyRates)


		for _, rate := range currencyRates {

			if rate.CurrencyCodeL == currency {
				dateRateMap.Store(fmt.Sprintf("%s/%s", currency, nbuDateFormatString), rate.Amount)
				//dateRateMap[fmt.Sprintf("%s/%s", currency, nbuDateFormatString)] = rate.Amount
				return rate.Amount, nil
			}
		}
	}
	// if no currency rate found
	return float64(0), errors.New("No currency found")
}
