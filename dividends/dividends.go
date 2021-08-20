package dividends

import (
	"fmt"
	"github.com/neshcheretal/ffparser/jsonreport"
	"github.com/neshcheretal/ffparser/nbu"
	"os"
	"time"
)

// Type for broker transaction (separate stock buy or sell orders)
type CashFlow struct {
	Date     time.Time
	Account  string
	Amount   float64
	Currency string
	UahRate  float64
	TypeId   string
	Comment  string
}

// type CashFlowResult struct {
//     Flows []CashFlow
//     Error error
// }

//Receive pased execl report
func CashFlowParser(jsonParsedFlows []jsonreport.JsonCashFlow) ([]CashFlow, error) {
	resultFlows := make([]CashFlow, 0)
	//dateRateMap := make(map[string]float64)
	dateLayout := "2006-01-02"
	for _, flow := range jsonParsedFlows {
		flowDate, err := time.Parse(dateLayout, flow.Date)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		flowAccount := flow.Account
		flowAmount := flow.Amount
		flowCurrency := flow.Currency

		flowUsdPriceUah, err := nbu.GetConversionRates(flowDate, flowCurrency)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		flowTypeId := flow.Type_id
		flowComment := flow.Comment
		if flowTypeId == "dividend" {
			newFlow := CashFlow{
				flowDate,
				flowAccount,
				flowAmount,
				flowCurrency,
				flowUsdPriceUah,
				flowTypeId,
				flowComment,
			}
			resultFlows = append(resultFlows, newFlow)
		}

	}
	return resultFlows, nil
}

func FilterDateFlows(flows []CashFlow, startDate time.Time, endDate time.Time) []CashFlow {
	filteredFlows := make([]CashFlow, 0)
	if startDate.IsZero() && endDate.IsZero() {
		filteredFlows = flows
	} else if !startDate.IsZero() && endDate.IsZero() {
		for _, flow := range flows {
			if flow.Date.After(startDate) {
				filteredFlows = append(filteredFlows, flow)
			}
		}
	} else if startDate.IsZero() && !endDate.IsZero() {
		for _, flow := range flows {
			if flow.Date.Before(endDate) {
				filteredFlows = append(filteredFlows, flow)
			}
		}
	} else {
		for _, flow := range flows {
			if flow.Date.After(startDate) && flow.Date.Before(endDate) {
				filteredFlows = append(filteredFlows, flow)
			}
		}
	}
	return filteredFlows
}

func DividendsPreparationWrapper(flows []jsonreport.JsonCashFlow, startDate time.Time, endDate time.Time, ch chan []CashFlow) {
	allDividendFlows, err := CashFlowParser(flows)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dateFilteredDividendFlows := FilterDateFlows(allDividendFlows, startDate, endDate)
	ch <- dateFilteredDividendFlows

}
