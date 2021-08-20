package dividends

import (
	"github.com/google/go-cmp/cmp"
	"github.com/neshcheretal/ffparser/jsonreport"
	"github.com/neshcheretal/ffparser/nbu"
	"github.com/neshcheretal/ffparser/utils/mocks"
	"testing"
	"time"
)

func TestCashFlowParser(t *testing.T) {
	nbu.Client = &mocks.MockClient{
		MockGet: mocks.NbuMainMock,
	}

	testCases := []struct {
		Name      string
		TestFlows []jsonreport.JsonCashFlow
		Expected  []CashFlow
	}{
		{
			Name:      "one_dividend",
			TestFlows: []jsonreport.JsonCashFlow{jsonreport.JsonCashFlow{"2020-09-03", "trade", 10.0, "USD", "dividend", "Dividends from kind people"}},
			Expected: []CashFlow{
				CashFlow{time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "trade", 10.0, "USD", 27.6428, "dividend", "Dividends from kind people"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got, _ := CashFlowParser(testCase.TestFlows)
			if !cmp.Equal(got, testCase.Expected) {
				t.Errorf("got %v, expected %v", got, testCase.Expected)
			}
		})
	}
}

func TestFilterDateFlows(t *testing.T) {
	testCashFlowList := []CashFlow{
		CashFlow{time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "trade", 10.0, "USD", 27.6428, "dividend", "Dividends from kind people"},
	}
	testCases := []struct {
		Name      string
		TestFlows []CashFlow
		StartDate time.Time
		EndDate   time.Time
		Expected  []CashFlow
	}{
		{
			Name:      "both_dates_dividend_inside",
			TestFlows: testCashFlowList,
			StartDate: time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
			Expected:  testCashFlowList,
		},
		{
			Name:      "both_dates_dividend_outside",
			TestFlows: testCashFlowList,
			StartDate: time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2020, time.September, 6, 0, 0, 0, 0, time.UTC),
			Expected:  []CashFlow{},
		},
		{
			Name:      "start_date_dividend_inside",
			TestFlows: testCashFlowList,
			StartDate: time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Time{},
			Expected:  testCashFlowList,
		},
		{
			Name:      "start_date_dividend_outside",
			TestFlows: testCashFlowList,
			StartDate: time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Time{},
			Expected:  []CashFlow{},
		},
		{
			Name:      "end_date_dividend_inside",
			TestFlows: testCashFlowList,
			StartDate: time.Time{},
			EndDate:   time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
			Expected:  testCashFlowList,
		},
		{
			Name:      "end_date_dividend_outside",
			TestFlows: testCashFlowList,
			StartDate: time.Time{},
			EndDate:   time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			Expected:  []CashFlow{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got := FilterDateFlows(testCase.TestFlows, testCase.StartDate, testCase.EndDate)
			if !cmp.Equal(got, testCase.Expected) {
				t.Errorf("got %v, expected %v", got, testCase.Expected)
			}
		})
	}
}

func TestDividendsPreparationWrapper(t *testing.T) {
	nbu.Client = &mocks.MockClient{
		MockGet: mocks.NbuMainMock,
	}
	testCh := make(chan []CashFlow, 1)

	testCases := []struct {
		Name      string
		TestFlows []jsonreport.JsonCashFlow
		StartDate time.Time
		EndDate   time.Time
		Expected  []CashFlow
	}{
		{
			Name:      "one_dividend",
			TestFlows: []jsonreport.JsonCashFlow{jsonreport.JsonCashFlow{"2020-09-03", "trade", 10.0, "USD", "dividend", "Dividends from kind people"}},
			StartDate: time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
			Expected: []CashFlow{
				CashFlow{time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "trade", 10.0, "USD", 27.6428, "dividend", "Dividends from kind people"},
			},
		},
		{
			Name:      "one_dividend_outside",
			TestFlows: []jsonreport.JsonCashFlow{jsonreport.JsonCashFlow{"2020-09-03", "trade", 10.0, "USD", "dividend", "Dividends from kind people"}},
			StartDate: time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2020, time.September, 6, 0, 0, 0, 0, time.UTC),
			Expected:  []CashFlow{},
		},
		{
			Name:      "start_date_one_dividend_inside",
			TestFlows: []jsonreport.JsonCashFlow{jsonreport.JsonCashFlow{"2020-09-03", "trade", 10.0, "USD", "dividend", "Dividends from kind people"}},
			StartDate: time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Time{},
			Expected: []CashFlow{
				CashFlow{time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "trade", 10.0, "USD", 27.6428, "dividend", "Dividends from kind people"},
			},
		},
		{
			Name:      "start_date_one_dividend_outside",
			TestFlows: []jsonreport.JsonCashFlow{jsonreport.JsonCashFlow{"2020-09-03", "trade", 10.0, "USD", "dividend", "Dividends from kind people"}},
			StartDate: time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Time{},
			Expected:  []CashFlow{},
		},
		{
			Name:      "end_date_one_dividend_inside",
			TestFlows: []jsonreport.JsonCashFlow{jsonreport.JsonCashFlow{"2020-09-03", "trade", 10.0, "USD", "dividend", "Dividends from kind people"}},
			StartDate: time.Time{},
			EndDate:   time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
			Expected: []CashFlow{
				CashFlow{time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC), "trade", 10.0, "USD", 27.6428, "dividend", "Dividends from kind people"},
			},
		},
		{
			Name:      "end_date_one_dividend_outside",
			TestFlows: []jsonreport.JsonCashFlow{jsonreport.JsonCashFlow{"2020-09-03", "trade", 10.0, "USD", "dividend", "Dividends from kind people"}},
			StartDate: time.Time{},
			EndDate:   time.Date(2020, time.September, 2, 0, 0, 0, 0, time.UTC),
			Expected:  []CashFlow{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			DividendsPreparationWrapper(testCase.TestFlows, testCase.StartDate, testCase.EndDate, testCh)
			got := <-testCh
			if !cmp.Equal(got, testCase.Expected) {
				t.Errorf("got %v, expected %v", got, testCase.Expected)
			}
		})
	}
}
