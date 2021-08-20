package trading

import (
	"github.com/google/go-cmp/cmp"
	"github.com/neshcheretal/ffparser/jsonreport"
	"testing"
	"time"
)

func TestSplitsStockParser(t *testing.T) {
	testCases := []struct {
		Name         string
		TestSplits   []jsonreport.JsonSecurityIN
		Expected     []SecurityIN
		ExpectsError bool
	}{
		{
			Name: "Rename/same count",
			TestSplits: []jsonreport.JsonSecurityIN{
				jsonreport.JsonSecurityIN{
					"-1",
					"TEST.US",
					"split",
					"2020-09-05 15:00:01",
					"Reorg New Symbol: 1 TESTA (AAA1234B1234) -> 1 TESTB (12345C123) SD 05.09.2020",
				},
			},
			Expected: []SecurityIN{
				SecurityIN{-1, "TESTA", 1, "TESTB", 1, time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC)},
			},
			ExpectsError: false,
		},
		{
			Name: "Same name/different count",
			TestSplits: []jsonreport.JsonSecurityIN{
				jsonreport.JsonSecurityIN{
					"-1",
					"TEST.US",
					"split",
					"2020-09-05 15:00:01",
					"Reorg New Symbol: 1 TEST (AAA1234B1234) -> 2 TEST (12345C123) SD 05.09.2020",
				},
			},
			Expected: []SecurityIN{
				SecurityIN{-1, "TEST", 1, "TEST", 2, time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC)},
			},
			ExpectsError: false,
		},
		{
			Name: "Wrong Quantity",
			TestSplits: []jsonreport.JsonSecurityIN{
				jsonreport.JsonSecurityIN{"Wrong_here", "TEST.US", "split", "2020-09-05 15:00:01", "Reorg New Symbol: 1 TEST (AAA1234B1234) -> 2 TEST (12345C123) SD 05.09.2020"},
			},
			Expected:     nil,
			ExpectsError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got, err := SplitsStockParser(testCase.TestSplits)
			if testCase.ExpectsError {
				if err == nil {
					t.Errorf("Error is expected, got %v", err)
				}
			}
			if !cmp.Equal(got, testCase.Expected) {
				t.Errorf("got %v, expected %v", got, testCase.Expected)
			}
		})
	}
}

func TestStockOrderSplitReEvaluation(t *testing.T) {
	testCases := []struct {
		Name       string
		TestOrders []Order
		TestSplits []SecurityIN
		Expected   []Order
	}{
		{
			Name: "One split/Rename/same count",
			TestOrders: []Order{
				Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TESTA", "buy", 1, 100.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TESTA", "buy", 1, 100.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 6, 16, 30, 0, 0, time.UTC), "TESTB", "sell", 2, 100.0, 27.0, "USD", 2.0},
			},
			TestSplits: []SecurityIN{
				SecurityIN{1, "TESTA", 1, "TESTB", 1, time.Date(2020, time.September, 5, 15, 00, 0, 0, time.UTC)},
			},
			Expected: []Order{
				Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TESTB", "buy", 1, 100.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TESTB", "buy", 1, 100.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 6, 16, 30, 0, 0, time.UTC), "TESTB", "sell", 2, 100.0, 27.0, "USD", 2.0},
			},
		},
		{
			Name: "One split/Same Name/New count",
			TestOrders: []Order{
				Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TESTA", "buy", 1, 100.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TESTA", "buy", 1, 100.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 6, 16, 30, 0, 0, time.UTC), "TESTA", "sell", 4, 50.0, 27.0, "USD", 2.0},
			},
			TestSplits: []SecurityIN{
				SecurityIN{1, "TESTA", 1, "TESTA", 2, time.Date(2020, time.September, 5, 15, 00, 0, 0, time.UTC)},
			},
			Expected: []Order{
				Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TESTA", "buy", 2, 50.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TESTA", "buy", 2, 50.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 6, 16, 30, 0, 0, time.UTC), "TESTA", "sell", 4, 50.0, 27.0, "USD", 2.0},
			},
		},
		{
			Name: "Two split/New name/Same count/Same Name/New count",
			TestOrders: []Order{
				Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TESTA", "buy", 1, 100.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TESTA", "buy", 1, 100.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 6, 16, 30, 0, 0, time.UTC), "TESTB", "buy", 1, 100.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 8, 16, 30, 0, 0, time.UTC), "TESTB", "sell", 6, 50.0, 27.0, "USD", 2.0},
			},
			TestSplits: []SecurityIN{
				SecurityIN{1, "TESTA", 1, "TESTB", 1, time.Date(2020, time.September, 5, 15, 00, 0, 0, time.UTC)},
				SecurityIN{1, "TESTB", 1, "TESTB", 2, time.Date(2020, time.September, 7, 15, 00, 0, 0, time.UTC)},
			},
			Expected: []Order{
				Order{time.Date(2020, time.September, 3, 16, 30, 0, 0, time.UTC), "TESTB", "buy", 2, 50.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 4, 16, 30, 0, 0, time.UTC), "TESTB", "buy", 2, 50.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 6, 16, 30, 0, 0, time.UTC), "TESTB", "buy", 2, 50.0, 27.0, "USD", 2.0},
				Order{time.Date(2020, time.September, 8, 16, 30, 0, 0, time.UTC), "TESTB", "sell", 6, 50.0, 27.0, "USD", 2.0},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			stockOrderSplitReEvaluation(testCase.TestSplits, testCase.TestOrders)
			// function modifies testCase.TestOrders
			if !cmp.Equal(testCase.TestOrders, testCase.Expected) {
				t.Errorf("got %v, expected %v", testCase.TestOrders, testCase.Expected)
			}
		})
	}
}
