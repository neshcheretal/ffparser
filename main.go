package main

import (
	"fmt"
	"github.com/neshcheretal/ffparser/dividends"
	"github.com/neshcheretal/ffparser/excel"
	"github.com/neshcheretal/ffparser/jsonreport"
	"github.com/neshcheretal/ffparser/trading"
	"github.com/neshcheretal/ffparser/validator"
)

func main() {
	// Get user input and convert it to a proper format
	inputValues, err := validator.InputParse()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse the broker report JSON file
	jsonParsedReport, err := jsonreport.ParseJsonReport(inputValues.FileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Define channels for concurrent trade and divident calculation
	divCh := make(chan []dividends.CashFlow)
	tradeCh := make(chan map[string][]trading.Trade)

	go dividends.DividendsPreparationWrapper(jsonParsedReport.Cash_flows.Detailed, inputValues.StartDate, inputValues.EndDate, divCh)
	go trading.StockTradePreparationWrapper(jsonParsedReport, inputValues.StartDate, inputValues.EndDate, tradeCh)

	// Pass channels with results of calculation to excel report prep function
	tax_report := excel.ExelCreate(tradeCh, divCh, inputValues.ReportLang)
	if err := tax_report.SaveAs(inputValues.OutputFile); err != nil {
		fmt.Println(err)
	}
}
