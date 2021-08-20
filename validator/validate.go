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
	startDate := flag.String("start", "", "A start day for trade calculation in YYYY-MM-DD format. Still requires a full report as ned access to trade open date orders (Optional)")
	endDate := flag.String("end", "", "A start day for trade calculation in YYYY-MM-DD format. Still requires a full report as ned access to trade open date orders (Optional)")
	outputFile := flag.String("output", "tax_calculation.xlsx", "Name of output xlsx file for calculation report")
	flag.Parse()

	// validate input
	if *fileName == "" || *fileName == " " {
		return InputArgs{}, errors.New("No broker report file was provided")
	}

	dateLayout := "2006-01-02"
	var err error
	var startDateFormat time.Time
	if *startDate != "" {
		startDateFormat, err = time.Parse(dateLayout, *startDate)
		if err != nil {
			return InputArgs{}, err
		}
	}

	var endDateFormat time.Time
	if *endDate != "" {
		endDateFormat, err = time.Parse(dateLayout, *endDate)
		if err != nil {
			return InputArgs{}, err
		}
	}

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
