// package excel
//
// import (
// 	"github.com/google/go-cmp/cmp"
// 	"github.com/neshcheretal/ffparser/dividends"
// 	"github.com/neshcheretal/ffparser/language"
// 	"github.com/neshcheretal/ffparser/trading"
// 	"github.com/xuri/excelize/v2"
// 	"testing"
// 	"time"
// )
//
// func exelParse(filename string, sheet string) ([][]string, error) {
// 	f, err := excelize.OpenFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Get all the rows in the Sheet1.
// 	rows, err := f.GetRows(sheet)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return rows, nil
// }
//
// func TestExelCreate(t *testing.T) {
// 	testTradeMap := map[string][]trading.Trade{
// 		"TEST": []trading.Trade{
// 			trading.Trade{
// 				time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC),
// 				time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC),
// 				100.0,
// 				27.6428,
// 				101.0,
// 				27.6908,
// 				1,
// 				2.0,
// 				2.0,
// 			},
// 		},
// 	}
// 	testDividendList := []dividends.CashFlow{
// 		dividends.CashFlow{time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "trade", 10.0, "USD", 27.6428, "dividend", "Dividends from kind people"},
// 	}
// 	tch := make(chan map[string][]trading.Trade, 1)
// 	dch := make(chan []dividends.CashFlow, 1)
//
// 	testCases := []struct {
// 		Name          string
// 		ReferenceFile string
// 		Lang          language.Language
// 		TradesMap     map[string][]trading.Trade
// 		DividendList  []dividends.CashFlow
// 	}{
// 		{
// 			Name:          "One trade report RU",
// 			ReferenceFile: "testreports/testru.xlsx",
// 			Lang: language.Language{
// 				language.RuLang,
// 			},
// 			TradesMap:    testTradeMap,
// 			DividendList: testDividendList,
// 		},
// 		{
// 			Name:          "One trade report UA",
// 			ReferenceFile: "testreports/testua.xlsx",
// 			Lang: language.Language{
// 				language.UaLang,
// 			},
// 			TradesMap:    testTradeMap,
// 			DividendList: testDividendList,
// 		},
// 		{
// 			Name:          "One trade report EN",
// 			ReferenceFile: "testreports/testen.xlsx",
// 			Lang: language.Language{
// 				language.EnLang,
// 			},
// 			TradesMap:    testTradeMap,
// 			DividendList: testDividendList,
// 		},
// 	}
//
// 	for _, testCase := range testCases {
// 		t.Run(testCase.Name, func(t *testing.T) {
// 			// Get a value to compare from a verified file
// 			sheet1Name := testCase.Lang.Dictionary["Sheet1"]
// 			sheet2Name := testCase.Lang.Dictionary["Sheet2"]
// 			sheet3Name := testCase.Lang.Dictionary["Sheet3"]
// 			referenceFileRowsSheet1, err := exelParse(testCase.ReferenceFile, sheet1Name)
// 			if err != nil {
// 				t.Errorf("Error with reference file, got %v", err)
// 			}
// 			referenceFileRowsSheet2, err := exelParse(testCase.ReferenceFile, sheet2Name)
// 			if err != nil {
// 				t.Errorf("Error with reference file, got %v", err)
// 			}
// 			referenceFileRowsSheet3, err := exelParse(testCase.ReferenceFile, sheet3Name)
// 			if err != nil {
// 				t.Errorf("Error with reference file, got %v", err)
// 			}
//
// 			// Get a function result
// 			tch <- testCase.TradesMap
// 			dch <- testCase.DividendList
// 			got := ExelCreate(tch, dch, testCase.Lang)
// 			gotRowsSheet1, err := got.GetRows(sheet1Name)
// 			if err != nil {
// 				t.Errorf("Error with function output rows, got %v", err)
// 			}
//
// 			if !cmp.Equal(gotRowsSheet1, referenceFileRowsSheet1) {
// 				t.Errorf("got %v, expected %v", gotRowsSheet1, referenceFileRowsSheet1)
// 			}
//
// 			gotRowsSheet2, err := got.GetRows(sheet2Name)
// 			if err != nil {
// 				t.Errorf("Error with function output rows, got %v", err)
// 			}
//
// 			if !cmp.Equal(gotRowsSheet2, referenceFileRowsSheet2) {
// 				t.Errorf("got %v, expected %v", gotRowsSheet2, referenceFileRowsSheet2)
// 			}
//
// 			gotRowsSheet3, err := got.GetRows(sheet3Name)
// 			if err != nil {
// 				t.Errorf("Error with function output rows, got %v", err)
// 			}
//
// 			if !cmp.Equal(gotRowsSheet3, referenceFileRowsSheet3) {
// 				t.Errorf("got %v, expected %v", gotRowsSheet3, referenceFileRowsSheet3)
// 			}
// 		})
// 	}
// }
