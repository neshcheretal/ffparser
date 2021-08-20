package mocks

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

// Custom type that allows setting the func that our Mock Do func will run instead
type MockGetType func(url string) (resp *http.Response, err error)

// MockClient is the mock client
type MockClient struct {
	MockGet MockGetType
}

// Overriding what the Get function should "do" in our MockClient
func (m *MockClient) Get(url string) (resp *http.Response, err error) {
	return m.MockGet(url)
}

func NbuMainMock(url string) (*http.Response, error) {
	nbuJsonResponse1 := `[{
        "StartDate":"09.03.2020","TimeSign":"0000","CurrencyCode":"840","CurrencyCodeL":"USD","Units":1,"Amount":27.6428
    }]`
	nbuJsonResponse2 := `[{
        "StartDate":"09.04.2020","TimeSign":"0000","CurrencyCode":"840","CurrencyCodeL":"USD","Units":1,"Amount":27.6908
    }]`
	if url == "https://bank.gov.ua/NBU_Exchange/exchange?json&date=03092020" {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(nbuJsonResponse1))),
		}, nil
	} else if url == "https://bank.gov.ua/NBU_Exchange/exchange?json&date=04092020" {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(nbuJsonResponse2))),
		}, nil
	}
	return &http.Response{}, errors.New("Mock for this date is not implemented")
}
