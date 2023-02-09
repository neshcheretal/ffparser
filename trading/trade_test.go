package trading

import (
	"github.com/google/go-cmp/cmp"
	"github.com/neshcheretal/ffparser/jsonreport"
	"github.com/neshcheretal/ffparser/nbu"
	"github.com/neshcheretal/ffparser/utils/mocks"
	"testing"
	"time"
)

func TestAllTradeListPrepare(t *testing.T) {
	testCases := []struct {
		Name       string
		TestOrders StockOrders
		Expected   []Trade
	}{
		{
			Name: "1 buy/1sell",
			TestOrders: StockOrders{
				[]Order{Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "TEST", "buy", 1, 100.0, 27.6428, "USD", 2.0}},
				1,
				[]Order{Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC), "TEST", "sell", 1, 101.0, 27.6908, "USD", 2.}},
				1,
			},
			Expected: []Trade{
				Trade{
					time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
					100.0,
					27.6428,
					101.0,
					27.6908,
					1,
					2.0,
					2.0,
				},
			},
		},
		{
			Name: "1 buy/1buy/1sell",
			TestOrders: StockOrders{
				[]Order{
					Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "TEST", "buy", 1, 100.0, 27.6428, "USD", 2.0},
					Order{time.Date(2020, time.September, 3, 16, 35, 0, 0, time.UTC), time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "TEST", "buy", 1, 100.5, 27.6428, "USD", 2.0},
				},
				1,
				[]Order{Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC),time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC), "TEST", "sell", 1, 101.0, 27.6908, "USD", 2.}},
				1,
			},
			Expected: []Trade{
				Trade{
					time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
					100.0,
					27.6428,
					101.0,
					27.6908,
					1,
					2.0,
					2.0,
				},
			},
		},
		{
			Name: "1buy/1buy/2sell",
			TestOrders: StockOrders{
				[]Order{
					Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "TEST", "buy", 1, 100.0, 27.6428, "USD", 2.0},
					Order{time.Date(2020, time.September, 3, 16, 35, 0, 0, time.UTC), time.Date(2020, time.September, 3, 16, 35, 0, 0, time.UTC), "TEST", "buy", 1, 100.5, 27.6428, "USD", 2.0},
				},
				1,
				[]Order{Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC), "TEST", "sell", 2, 101.0, 27.6908, "USD", 2.}},
				1,
			},
			Expected: []Trade{
				Trade{
					time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
					100.0,
					27.6428,
					101.0,
					27.6908,
					1,
					2.0,
					1.0, // sell comission 2.0 divided betwen two closed trades
				},
				Trade{
					time.Date(2020, time.September, 3, 16, 35, 0, 0, time.UTC),
					time.Date(2020, time.September, 3, 16, 35, 0, 0, time.UTC),
					time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
					100.5,
					27.6428,
					101.0,
					27.6908,
					1,
					2.0,
					1.0, // sell comission 2.0 divided betwen two closed trades
				},
			},
		},
		{
			Name: "2buy/1sell/3buy/2sell",
			TestOrders: StockOrders{
				[]Order{
					// 2 buy
					Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "TEST", "buy", 2, 100.0, 27.6428, "USD", 2.0},
					// 3 buy
					Order{time.Date(2020, time.September, 7, 16, 30, 0, 0, time.UTC), time.Date(2020, time.September, 7, 16, 30, 0, 0, time.UTC), "TEST", "buy", 3, 102.0, 27.7325, "USD", 2.1},
				},
				1,
				[]Order{
					// 1sell
					Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC), "TEST", "sell", 1, 101.0, 27.6908, "USD", 2.0},
					// 2sell
					Order{time.Date(2020, time.September, 8, 16, 30, 0, 0, time.UTC), time.Date(2020, time.September, 8, 16, 30, 0, 0, time.UTC), "TEST", "sell", 2, 104.0, 27.7509, "USD", 2.0},
				},
				1,
			},
			Expected: []Trade{
				Trade{
					time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
					100.0,
					27.6428,
					101.0,
					27.6908,
					1,
					1.0, // buy comission 2.0 from the first buy is divided betwen two closed trades
					2.0,
				},
				Trade{
					time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2020, time.September, 8, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 8, 16, 30, 0, 0, time.UTC),
					100.0,
					27.6428,
					104.0,
					27.7509,
					1,
					1.0, // buy comission 2.0 from the first buy is divided betwen two closed trades
					1.0, // sell comission 2.0 from the second sell is divided betwen two closed trades
				},
				Trade{
					time.Date(2020, time.September, 7, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 7, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 8, 16, 30, 0, 0, time.UTC),
					time.Date(2020, time.September, 8, 16, 30, 0, 0, time.UTC),
					102.0,
					27.7325,
					104.0,
					27.7509,
					1,
					2.1 / float64(3), // buy comission 2.1 from the second buy is divided on 3 as we sell one third of amount, the remaining (2.0/float64(3))*2 is left for possible next trades
					1.0,              // sell comission 2.0 from the second sell is divided betwen two closed trades
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got := AllTradeListPrepare(testCase.TestOrders)
			if !cmp.Equal(got, testCase.Expected) {
				t.Errorf("got %v, expected %v", got, testCase.Expected)
			}
		})
	}
}

func TestFilterDateTrades(t *testing.T) {

	TestTradeList := []Trade{
		Trade{
			time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC),
			time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC),
			time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
			100.0,
			27.6428,
			101.0,
			27.6908,
			1,
			2.0,
			1.0, // sell comission 2.0 divided betwen two closed trades
		},
	}

	testCases := []struct {
		Name       string
		TestTrades []Trade
		StartDate  time.Time
		EndDate    time.Time
		Expected   []Trade
	}{
		{
			Name:       "both_dates_trade_inside",
			TestTrades: TestTradeList,
			StartDate:  time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC),
			Expected:   TestTradeList,
		},
		{
			Name:       "both_dates_trade_outside",
			TestTrades: TestTradeList,
			StartDate:  time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2020, time.September, 6, 0, 0, 0, 0, time.UTC),
			Expected:   []Trade{},
		},
		{
			Name:       "start_date_trade_inside",
			TestTrades: TestTradeList,
			StartDate:  time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Time{},
			Expected:   TestTradeList,
		},
		{
			Name:       "start_date_trade_outside",
			TestTrades: TestTradeList,
			StartDate:  time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Time{},
			Expected:   []Trade{},
		},
		{
			Name:       "end_date_trade_inside",
			TestTrades: TestTradeList,
			StartDate:  time.Time{},
			EndDate:    time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC),
			Expected:   TestTradeList,
		},
		{
			Name:       "end_date_trade_outside",
			TestTrades: TestTradeList,
			StartDate:  time.Time{},
			EndDate:    time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			Expected:   []Trade{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got := FilterDateTrades(testCase.TestTrades, testCase.StartDate, testCase.EndDate)
			if !cmp.Equal(got, testCase.Expected) {
				t.Errorf("got %v, expected %v", got, testCase.Expected)
			}
		})
	}
}

func TestStockTradePreparationWrapper(t *testing.T) {
	nbu.Client = &mocks.MockClient{
		MockGet: mocks.NbuMainMock,
	}

	testCh := make(chan map[string][]Trade, 1)

	testOrderList := []jsonreport.JsonOrder{jsonreport.JsonOrder{"2020-09-03 16:30:00", "2020-09-03", "TEST.US", "buy", 1, 100.0, "USD", 2.0, "USD"}, jsonreport.JsonOrder{"2020-09-04 16:30:00", "2020-09-04", "TEST.US", "sell", 1, 101.0, "USD", 2.0, "USD"}}
	testJsonReport := jsonreport.JsonReport{
		jsonreport.JsonDetailedTradeReport{testOrderList},
		jsonreport.JsonDetailedCashReport{},
		[]jsonreport.JsonSecurityIN{},
	}
	testTradeMap := map[string][]Trade{
		"TEST": []Trade{
			Trade{
				time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC),
				time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC),
				time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
				100.0,
				27.6428,
				101.0,
				27.6908,
				1,
				2.0,
				2.0,
			},
		},
	}
	testEmptyTradeMap := make(map[string][]Trade)

	testCases := []struct {
		Name       string
		TestReport jsonreport.JsonReport
		StartDate  time.Time
		EndDate    time.Time
		Expected   map[string][]Trade
	}{
		{
			Name:       "one_trade_inside",
			TestReport: testJsonReport,
			StartDate:  time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC),
			Expected:   testTradeMap,
		},
		{
			Name:       "one_trade_outside",
			TestReport: testJsonReport,
			StartDate:  time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2020, time.September, 6, 0, 0, 0, 0, time.UTC),
			Expected:   testEmptyTradeMap,
		},
		{
			Name:       "start_date_one_trade_inside",
			TestReport: testJsonReport,
			StartDate:  time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Time{},
			Expected:   testTradeMap,
		},
		{
			Name:       "start_date_one_trade_outside",
			TestReport: testJsonReport,
			StartDate:  time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Time{},
			Expected:   testEmptyTradeMap,
		},
		{
			Name:       "end_date_one_trade_inside",
			TestReport: testJsonReport,
			StartDate:  time.Time{},
			EndDate:    time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC),
			Expected:   testTradeMap,
		},
		{
			Name:       "end_date_one_trade_outside",
			TestReport: testJsonReport,
			StartDate:  time.Time{},
			EndDate:    time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			Expected:   testEmptyTradeMap,
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.Name, func(t *testing.T) {
			StockTradePreparationWrapper(testCase.TestReport, testCase.StartDate, testCase.EndDate, testCh)
			got := <-testCh
			if !cmp.Equal(got, testCase.Expected) {
				t.Errorf("got %v, expected %v", got, testCase.Expected)
			}
		})
	}
}
