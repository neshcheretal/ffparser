package validator

import (
	"flag"
	"github.com/google/go-cmp/cmp"
	"github.com/neshcheretal/ffparser/language"
	"os"
	"testing"
	"time"
)

func TestInputParse(t *testing.T) {
	testCases := []struct {
		Name         string
		Flags        []string
		Expected     InputArgs
		ExpectsError bool
	}{
		{
			Name:  "File/year/RU/output",
			Flags: []string{" ", "-report=test.json", "-year=2020", "-lang=RU", "-output=test"},
			Expected: InputArgs{
				"test.json",
				language.Language{language.RuLang},
				time.Date(2019, time.December, 31, 23, 59, 59, 0, time.UTC),
				time.Date(2020, time.December, 31, 23, 59, 59, 0, time.UTC),
				"test.xlsx",
			},
			ExpectsError: false,
		},
		{
			Name:         "File/no year/EN/output",
			Flags:        []string{" ", "-report=test.json", "-lang=EN", "-output=test"},
			Expected:     InputArgs{},
			ExpectsError: true,
		},
		{
			Name:  "File/year/No lang/output",
			Flags: []string{" ", "-report=test.json", "-year=2020", "-output=test"},
			Expected: InputArgs{
				"test.json",
				language.Language{language.EnLang},
				time.Date(2019, time.December, 31, 23, 59, 59, 0, time.UTC),
				time.Date(2020, time.December, 31, 23, 59, 59, 0, time.UTC),
				"test.xlsx",
			},
			ExpectsError: false,
		},
		{
			Name:  "File/year/UA/no output",
			Flags: []string{" ", "-report=test.json", "-year=2020", "-lang=UA"},
			Expected: InputArgs{
				"test.json",
				language.Language{language.UaLang},
				time.Date(2019, time.December, 31, 23, 59, 59, 0, time.UTC),
				time.Date(2020, time.December, 31, 23, 59, 59, 0, time.UTC),
				"tax_calculation.xlsx",
			},
			ExpectsError: false,
		},
		{
			Name:         "File/year wrong format/RU/output",
			Flags:        []string{" ", "-report=test.json", "-year=202aaa", "-lang=RU", "-output=test"},
			Expected:     InputArgs{},
			ExpectsError: true,
		},
		{
			Name:         "No file/year/RU/output",
			Flags:        []string{" ", "-year=2020", "-lang=RU", "-output=test"},
			Expected:     InputArgs{},
			ExpectsError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// reset flags value before new tests
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			// pass new flags
			os.Args = testCase.Flags
			got, err := InputParse()
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
