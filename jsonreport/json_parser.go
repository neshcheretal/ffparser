package jsonreport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

type JsonCashFlow struct {
	Date     string
	Account  string
	Amount   float64
	Currency string
	Type_id  string
	Comment  string
}

type JsonOrder struct {
	Date                string
	Pay_d               string
	Instr_nm            string
	Operation           string
	Q                   int
	P                   float64
	Curr_c              string
	Commission          float64
	Commission_currency string
}

type JsonDetailedCashReport struct {
	Detailed []JsonCashFlow
}

type JsonDetailedTradeReport struct {
	Detailed []JsonOrder
}

type JsonSecurityIN struct {
	Quantity string
	Ticker   string
	Type     string
	Datetime string
	Comment  string
}

type JsonReport struct {
	Trades             JsonDetailedTradeReport
	Cash_flows         JsonDetailedCashReport
	Securities_in_outs []JsonSecurityIN
}

func ParseJsonReport(filename string) (JsonReport, error) {
	var jsonParsedReport JsonReport
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		return JsonReport{}, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	dec := json.NewDecoder(bytes.NewReader(byteValue))
	if err := dec.Decode(&jsonParsedReport); err != nil {
		return JsonReport{}, err
	}
	return jsonParsedReport, nil
}
