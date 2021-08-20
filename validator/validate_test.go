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
			Name:  "File/start/RU/output",
			Flags: []string{" ", "-report=test.json", "-start=2020-09-03", "-lang=RU", "-output=test"},
			Expected: InputArgs{
				"test.json",
				language.Language{language.RuLang},
				time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
				time.Time{},
				"test.xlsx",
			},
			ExpectsError: false,
		},
		{
			Name:  "File/end/EN/output",
			Flags: []string{" ", "-report=test.json", "-end=2020-09-03", "-lang=EN", "-output=test"},
			Expected: InputArgs{
				"test.json",
				language.Language{language.EnLang},
				time.Time{},
				time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
				"test.xlsx",
			},
			ExpectsError: false,
		},
		{
			Name:  "File/end/No lang/output",
			Flags: []string{" ", "-report=test.json", "-end=2020-09-03", "-output=test"},
			Expected: InputArgs{
				"test.json",
				language.Language{language.EnLang},
				time.Time{},
				time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
				"test.xlsx",
			},
			ExpectsError: false,
		},
		{
			Name:  "File/end/UA/no output",
			Flags: []string{" ", "-report=test.json", "-end=2020-09-03", "-lang=UA"},
			Expected: InputArgs{
				"test.json",
				language.Language{language.UaLang},
				time.Time{},
				time.Date(2020, time.September, 3, 0, 0, 0, 0, time.UTC),
				"tax_calculation.xlsx",
			},
			ExpectsError: false,
		},
		{
			Name:         "File/start wrong format/RU/output",
			Flags:        []string{" ", "-report=test.json", "-start=2020.09.03", "-lang=RU", "-output=test"},
			Expected:     InputArgs{},
			ExpectsError: true,
		},
		{
			Name:         "File/end wrong format/RU/output",
			Flags:        []string{" ", "-report=test.json", "-end=2020.09.03", "-lang=RU", "-output=test"},
			Expected:     InputArgs{},
			ExpectsError: true,
		},
		{
			Name:         "No file/end/RU/output",
			Flags:        []string{" ", "-end=2020-09-03", "-lang=RU", "-output=test"},
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
