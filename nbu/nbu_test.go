package nbu

import (
	"bytes"
	"errors"
	"github.com/neshcheretal/ffparser/utils/mocks"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestNbuApiCallSuccess(t *testing.T) {
	// build our response JSON
	jsonResponse := `[{
        "StartDate":"16.08.2021","TimeSign":"0000","CurrencyCode":"840","CurrencyCodeL":"USD","Units":1,"Amount":26.6931
    }]`

	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonResponse)))
	Client = &mocks.MockClient{
		MockGet: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}

	result, err := GetConversionRates(time.Date(2021, time.August, 16, 12, 0, 0, 0, time.UTC), "USD")
	if err != nil {
		t.Error("TestNbuApiCallSuccess failed.")
		return
	}

	if result != 26.6931 {
		t.Error("GetConversionRates return wrong value")
		return
	}
}

func TestNbuApiCallFail(t *testing.T) { // create a client that throws and returns an error
	Client = &mocks.MockClient{
		MockGet: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       nil,
			}, errors.New("Mock Error")
		},
	}
	_, err := GetConversionRates(time.Date(2021, time.August, 16, 12, 0, 0, 0, time.UTC), "MONEY")
	if err == nil {
		t.Error("TestNbuApiCallFail failed.")
		return
	}
}
