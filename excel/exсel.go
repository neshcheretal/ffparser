package excel

import (
	"fmt"
	"github.com/neshcheretal/ffparser/dividends"
	"github.com/neshcheretal/ffparser/language"
	"github.com/neshcheretal/ffparser/trading"
	"github.com/xuri/excelize/v2"
)

type ColumnFormat struct {
	Width float64
	Head  string
}

var tradeColumnFormat = map[string]ColumnFormat{
	"A": ColumnFormat{
		13.5,
		"ticker",
	},
	"B": ColumnFormat{
		15.0,
		"tradeNumber",
	},
	"C": ColumnFormat{
		16.0,
		"tradeOpenDate",
	},
	"D": ColumnFormat{
		12.0,
		"tradeOpenPrice",
	},
	"E": ColumnFormat{
		16,
		"tradeOpenComission",
	},
	"F": ColumnFormat{
		18,
		"tradeOpenCurrencyRate",
	},
	"G": ColumnFormat{
		16.0,
		"tradeCloseDate",
	},
	"H": ColumnFormat{
		12.0,
		"tradeClosePrice",
	},
	"I": ColumnFormat{
		16.0,
		"tradeCloseComission",
	},
	"J": ColumnFormat{
		18,
		"tradeCloseCurrencyRate",
	},
	"K": ColumnFormat{
		9.5,
		"tradeStockQuantity",
	},
	"L": ColumnFormat{
		21.0,
		"tradeCloseCostUAH",
	},
	"M": ColumnFormat{
		21.0,
		"tradeOpenCostUAH",
	},
	"N": ColumnFormat{
		17.0,
		"tradeProfitUAH",
	},
}

var dividendColumnFormat = map[string]ColumnFormat{
	"A": ColumnFormat{
		19,
		"dividendDate",
	},
	"B": ColumnFormat{
		19,
		"dividendAmount",
	},
	"C": ColumnFormat{
		19,
		"dividendUahRate",
	},
	"D": ColumnFormat{
		65.0,
		"dividendComment",
	},
	"E": ColumnFormat{
		19.0,
		"dividendUAH",
	},
}

var tradeTaxColumnFormat = map[string]ColumnFormat{
	"A": ColumnFormat{
		40,
		"tradeTotalUahProfit",
	},
	"B": ColumnFormat{
		18.0,
		"tradeIncomeTax",
	},
	"C": ColumnFormat{
		26.0,
		"tradeMilitaryTax",
	},
}

var dividendTaxColumnFormat = map[string]ColumnFormat{
	"A": ColumnFormat{
		40.0,
		"dividendTotalUahProfit",
	},
	"B": ColumnFormat{
		18.0,
		"dividendIncomeTax",
	},
	"C": ColumnFormat{
		26.0,
		"dividendMilitaryTax",
	},
}

func ExelCreate(tch <-chan map[string][]trading.Trade, dch <-chan []dividends.CashFlow, reportLang language.Language) *excelize.File {
	f := excelize.NewFile()
	sheet1Name := reportLang.Dictionary["Sheet1"]
	sheet2Name := reportLang.Dictionary["Sheet2"]
	sheet3Name := reportLang.Dictionary["Sheet3"]
	f.SetSheetName("Sheet1", sheet1Name)

	// Set table headers format
	style, err := f.NewStyle(`{"font":{"bold":true}}`)
	if err != nil {
		fmt.Printf("Failed to set font style: %v", err)
	}
	for k, v := range tradeColumnFormat {
		f.SetColWidth(sheet1Name, k, k, v.Width)
		f.SetCellValue(sheet1Name, k+"1", reportLang.Dictionary[v.Head])
		f.SetCellStyle(sheet1Name, k+"1", k+"1", style)
	}

	tradeIndex := 2 // start from 2 as 1 is for table headers
	// Calculate per trade profit
	for stock, tradelist := range <-tch {
		for i, trade := range tradelist {
			f.SetCellValue(sheet1Name, fmt.Sprintf("A%d", tradeIndex), stock)
			f.SetCellValue(sheet1Name, fmt.Sprintf("B%d", tradeIndex), fmt.Sprintf("trade %d", i+1))
			f.SetCellValue(sheet1Name, fmt.Sprintf("C%d", tradeIndex), trade.OpenDate)
			f.SetCellValue(sheet1Name, fmt.Sprintf("D%d", tradeIndex), trade.OpenPrice)
			f.SetCellValue(sheet1Name, fmt.Sprintf("E%d", tradeIndex), trade.OpenComission)
			f.SetCellValue(sheet1Name, fmt.Sprintf("F%d", tradeIndex), trade.OpenUahRate)
			f.SetCellValue(sheet1Name, fmt.Sprintf("G%d", tradeIndex), trade.CloseDate)
			f.SetCellValue(sheet1Name, fmt.Sprintf("H%d", tradeIndex), trade.ClosePrice)
			f.SetCellValue(sheet1Name, fmt.Sprintf("I%d", tradeIndex), trade.CloseComission)
			f.SetCellValue(sheet1Name, fmt.Sprintf("J%d", tradeIndex), trade.CloseUahRate)
			f.SetCellValue(sheet1Name, fmt.Sprintf("K%d", tradeIndex), trade.Quantity)
			f.SetCellFormula(sheet1Name, fmt.Sprintf("L%d", tradeIndex), fmt.Sprintf("ROUND(((H%[1]d*K%[1]d-I%[1]d)*J%[1]d); 2)", tradeIndex))
			f.SetCellFormula(sheet1Name, fmt.Sprintf("M%d", tradeIndex), fmt.Sprintf("ROUND(((D%[1]d*K%[1]d+E%[1]d)*F%[1]d); 2)", tradeIndex))
			f.SetCellFormula(sheet1Name, fmt.Sprintf("N%d", tradeIndex), fmt.Sprintf("L%[1]d-M%[1]d", tradeIndex))
			tradeIndex += 1
		}
	}
	lastTradeIndex := tradeIndex - 1

	// Calculate dividends
	f.NewSheet(sheet2Name)
	for k, v := range dividendColumnFormat {
		f.SetColWidth(sheet2Name, k, k, v.Width)
		f.SetCellValue(sheet2Name, k+"1", reportLang.Dictionary[v.Head])
		f.SetCellStyle(sheet2Name, k+"1", k+"1", style)
	}

	dividendIndex := 2 // start from 2 as 1 is for table headers
	for _, dividend := range <-dch {
		f.SetCellValue(sheet2Name, fmt.Sprintf("A%d", dividendIndex), fmt.Sprintf(dividend.Date.Format("2006-01-02")))
		f.SetCellValue(sheet2Name, fmt.Sprintf("B%d", dividendIndex), dividend.Amount)
		f.SetCellValue(sheet2Name, fmt.Sprintf("C%d", dividendIndex), dividend.UahRate)
		f.SetCellValue(sheet2Name, fmt.Sprintf("D%d", dividendIndex), dividend.Comment)
		f.SetCellFormula(sheet2Name, fmt.Sprintf("E%d", dividendIndex), fmt.Sprintf("ROUND((B%[1]d*C%[1]d); 2)", dividendIndex))
		dividendIndex += 1
	}
	lastDividendIndex := dividendIndex - 1

	// Calculate tax amount
	f.NewSheet(sheet3Name)
	for k, v := range tradeTaxColumnFormat {
		f.SetColWidth(sheet3Name, k, k, v.Width)
		f.SetCellValue(sheet3Name, k+"1", reportLang.Dictionary[v.Head])
		f.SetCellStyle(sheet3Name, k+"1", k+"1", style)
	}
	f.SetCellFormula(sheet3Name, "A2", fmt.Sprintf("SUM(%v!N2:%v!N%d)", sheet1Name, sheet1Name, lastTradeIndex))
	f.SetCellFormula(sheet3Name, "B2", "ROUND((A2/100)*18; 2)")
	f.SetCellFormula(sheet3Name, "C2", "ROUND((A2/100)*1.5; 2)")

	for k, v := range dividendTaxColumnFormat {
		f.SetCellValue(sheet3Name, k+"4", reportLang.Dictionary[v.Head])
		f.SetCellStyle(sheet3Name, k+"4", k+"4", style)
	}
	f.SetCellFormula(sheet3Name, "A5", fmt.Sprintf("SUM(%v!E2:%v!E%d)", sheet2Name, sheet2Name, lastDividendIndex))
	f.SetCellFormula(sheet3Name, "B5", "ROUND((A5/100)*9; 2)")
	f.SetCellFormula(sheet3Name, "C5", "ROUND((A5/100)*1.5; 2)")
	return f
}
