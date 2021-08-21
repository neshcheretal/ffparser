package validator

import (
	"errors"
	"flag"
	"fmt"
	"github.com/neshcheretal/ffparser/language"
	"strings"
	"time"
)

type InputArgs struct {
	FileName   string
	ReportLang language.Language
	StartDate  time.Time
	EndDate    time.Time
	OutputFile string
}

func InputParse() (InputArgs, error) {
	fileName := flag.String("report", "", "A path to a FF broker report JSON file.  (Required)")
	reportLang := flag.String("lang", "EN", "A report language, now supported languages are UA, RU, EN.  (Required)")
	reportYear := flag.String("year", "", "A year for which tax should be calculated in YYYY format. Still requires a full report as need access to trade open date orders (Required)")
	outputFile := flag.String("output", "tax_calculation.xlsx", "Name of output xlsx file for calculation report")
	flag.Parse()

	// validate input
	if *fileName == "" || *fileName == " " {
		return InputArgs{}, errors.New("No broker report file was provided")
	}

	var err error
	var startDateFormat time.Time
	var endDateFormat time.Time
	yearLayout := "2006-01-02 15:04:05"
	if *reportYear == "" {
		return InputArgs{}, errors.New("You have to provide a year for which tax will be calculated")
	} else {
		// get last time of the year
		endDateFormat, err = time.Parse(yearLayout, fmt.Sprintf("%v-12-31 23:59:59", *reportYear))
		if err != nil {
			return InputArgs{}, err
		}
	}

	// get last time of the previous year
	startDateFormat = endDateFormat.AddDate(-1, 0, 0)

	var languageSet language.Language
	switch *reportLang {
	case "RU":
		languageSet = language.Language{language.RuLang}
	case "UA":
		languageSet = language.Language{language.UaLang}
	case "EN":
		languageSet = language.Language{language.EnLang}
	default:
		return InputArgs{}, errors.New(fmt.Sprintf("%s is not supported\n", *reportLang))
	}

	var outfileName string
	if strings.HasSuffix(*outputFile, ".xlsx") {
		outfileName = *outputFile
	} else {
		outfileName = *outputFile + ".xlsx"
	}

	validatedValues := InputArgs{
		*fileName,
		languageSet,
		startDateFormat,
		endDateFormat,
		outfileName,
	}
	return validatedValues, nil
}
