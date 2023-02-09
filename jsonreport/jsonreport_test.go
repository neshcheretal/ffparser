package jsonreport

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseJsonReport(t *testing.T) {
	testCases := []struct {
		Name         string
		Filename     string
		Expected     JsonReport
		ExpectsError bool
	}{
		{
			Name:     "One trade report",
			Filename: "testreports/broker-report-1.json",
			Expected: JsonReport{
				JsonDetailedTradeReport{
					[]JsonOrder{JsonOrder{"2020-09-03 16:30:00", "2020-09-03", "TEST.US", "buy", 1, 100.0, "USD", 2.0, "USD"}, JsonOrder{"2020-09-04 16:30:00", "2020-09-04", "TEST.US", "sell", 1, 101.0, "USD", 2.0, "USD"}},
				},
				JsonDetailedCashReport{},
				nil,
			},
			ExpectsError: false,
		},
		{
			Name:     "One trade report/one cash flow",
			Filename: "testreports/broker-report-2.json",
			Expected: JsonReport{
				JsonDetailedTradeReport{
					[]JsonOrder{JsonOrder{"2020-09-03 16:30:00", "2020-09-03", "TEST.US", "buy", 1, 100.0, "USD", 2.0, "USD"}, JsonOrder{"2020-09-04 16:30:00", "2020-09-04", "TEST.US", "sell", 1, 101.0, "USD", 2.0, "USD"}},
				},
				JsonDetailedCashReport{
					[]JsonCashFlow{JsonCashFlow{"2020-09-03", "торговый", 10, "USD", "dividend", "Dividends from kind people"}},
				},
				nil,
			},
			ExpectsError: false,
		},
		{
			Name:     "One trade report/one cash flow/one split",
			Filename: "testreports/broker-report-3.json",
			Expected: JsonReport{
				JsonDetailedTradeReport{
					[]JsonOrder{JsonOrder{"2020-09-03 16:30:00", "2020-09-03", "TEST.US", "buy", 1, 100.0, "USD", 2.0, "USD"}, JsonOrder{"2020-09-04 16:30:00", "2020-09-04", "TEST.US", "sell", 1, 101.0, "USD", 2.0, "USD"}},
				},
				JsonDetailedCashReport{
					[]JsonCashFlow{JsonCashFlow{"2020-09-03", "торговый", 10, "USD", "dividend", "Dividends from kind people"}},
				},
				[]JsonSecurityIN{
					JsonSecurityIN{
						"-1",
						"TEST.US",
						"split",
						"2020-09-05 15:00:01",
						"Reorg New Symbol: 1 TEST (AAA1234B1234) -> 1 TESTA (12345C123) SD 05.09.2020",
					},
				},
			},
			ExpectsError: false,
		},
		{
			Name:         "Non existing file",
			Filename:     "testreports/nosuchfile.json",
			Expected:     JsonReport{},
			ExpectsError: true,
		},
		{
			Name:         "Malformed json file",
			Filename:     "testreports/broker-report-malformed.json",
			Expected:     JsonReport{},
			ExpectsError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got, err := ParseJsonReport(testCase.Filename)
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
