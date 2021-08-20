package trading

import (
	"github.com/google/go-cmp/cmp"
	"github.com/neshcheretal/ffparser/jsonreport"
	"github.com/neshcheretal/ffparser/nbu"
	"github.com/neshcheretal/ffparser/utils/mocks"
	"testing"
	"time"
)

func TestOrderListPreparator(t *testing.T) {
	// mock NBU API calls
	nbu.Client = &mocks.MockClient{
		MockGet: mocks.NbuMainMock,
	}

	testCases := []struct {
		Name       string
		TestTrades []jsonreport.JsonOrder
		Expected   []Order
	}{
		{
			Name:       "Buy and sell trades",
			TestTrades: []jsonreport.JsonOrder{jsonreport.JsonOrder{"2020-09-03 16:30:00", "TEST", "buy", 1, 100.0, "USD", 2.0, "USD"}, jsonreport.JsonOrder{"2020-09-04 16:30:00", "TEST", "sell", 1, 101.0, "USD", 2.0, "USD"}},
			Expected: []Order{
				Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TEST", "buy", 1, 100.0, 27.6428, "USD", 2.0},
				Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TEST", "sell", 1, 101.0, 27.6908, "USD", 2.0},
			},
		},
		{
			Name:       "Only buy",
			TestTrades: []jsonreport.JsonOrder{jsonreport.JsonOrder{"2020-09-03 16:30:00", "TEST", "buy", 1, 100.0, "USD", 2.0, "USD"}},
			Expected:   []Order{Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TEST", "buy", 1, 100, 27.6428, "USD", 2}},
		},
		{
			Name:       "Only sell(when ticker was renamed)",
			TestTrades: []jsonreport.JsonOrder{jsonreport.JsonOrder{"2020-09-04 16:30:00", "TEST", "sell", 1, 101.0, "USD", 2.0, "USD"}},
			Expected:   []Order{Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TEST", "sell", 1, 101, 27.6908, "USD", 2}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got, _ := OrderListPreparator(testCase.TestTrades)
			if !cmp.Equal(got, testCase.Expected) {
				t.Errorf("got %v, expected %v", got, testCase.Expected)
			}
		})
	}
}

func TestOrdersStockParser(t *testing.T) {

	testCases := []struct {
		Name       string
		TestTrades []Order
		Expected   map[string]StockOrders
	}{
		{
			Name: "Buy and sell trades",
			TestTrades: []Order{
				Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TEST", "buy", 1, 100.0, 27.6428, "USD", 2.0},
				Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TEST", "sell", 1, 101.0, 27.6908, "USD", 2.0},
			},
			Expected: map[string]StockOrders{
				"TEST": StockOrders{
					[]Order{Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TEST", "buy", 1, 100.0, 27.6428, "USD", 2.0}},
					1,
					[]Order{Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TEST", "sell", 1, 101.0, 27.6908, "USD", 2.0}},
					1,
				},
			},
		},
		{
			Name:       "Only buy",
			TestTrades: []Order{Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TEST", "buy", 1, 100, 27.6428, "USD", 2}},
			Expected: map[string]StockOrders{
				"TEST": StockOrders{
					[]Order{Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TEST", "buy", 1, 100, 27.6428, "USD", 2}},
					1,
					nil,
					0,
				},
			},
		},
		{
			Name:       "Only sell(when ticker was renamed)",
			TestTrades: []Order{Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TEST", "sell", 1, 101, 27.6908, "USD", 2}},
			Expected: map[string]StockOrders{
				"TEST": StockOrders{
					nil,
					0,
					[]Order{Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TEST", "sell", 1, 101, 27.6908, "USD", 2}},
					1,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got := OrdersStockParser(testCase.TestTrades)
			if !cmp.Equal(got, testCase.Expected) {
				t.Errorf("got %v, expected %v", got, testCase.Expected)
			}
		})
	}
}
