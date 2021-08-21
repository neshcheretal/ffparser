package trading

import (
	//	"fmt"
	"github.com/neshcheretal/ffparser/jsonreport"
	"regexp"
	"strconv"
	"time"
	//    "errors"
)

// Type for broker transaction (separate stock buy or sell orders)
type SecurityIN struct {
	Quantity  int
	OldTicker string
	OldCount  int
	NewTicker string
	NewCount  int
	Datetime  time.Time
}

func SplitsStockParser(jsonParsedSecurities []jsonreport.JsonSecurityIN) ([]SecurityIN, error) {
	dateLayout := "02.01.2006"
	resultList := make([]SecurityIN, 0)
	regexpObject := regexp.MustCompile(".*: ([1-9]*) ([A-Z]*) \\(([A-Z0-9]*)\\) -> ([1-9]*) ([A-Z]*) \\(([A-Z0-9]*)\\) SD ([0-3][0-9]\\.[0-1][0-9]\\.[0-9]*)")
	for _, securities := range jsonParsedSecurities {
		if securities.Type != "split" {
			continue
		} else {

			splitAmount, err := strconv.Atoi(securities.Quantity)
			if err != nil {
				return nil, err
			}
			result := regexpObject.FindStringSubmatch(securities.Comment)
			splitOldCountString := result[1]
			splitOldCount, err := strconv.Atoi(splitOldCountString)
			if err != nil {
				return nil, err
			}
			splitOldTicker := result[2]
			splitNewCountString := result[4]
			splitNewCount, err := strconv.Atoi(splitNewCountString)
			if err != nil {
				return nil, err
			}
			splitNewTicker := result[5]
			splitDateString := result[7]
			splitDate, err := time.Parse(dateLayout, splitDateString)
			if err != nil {
				return nil, err
			}

			resultList = append(resultList, SecurityIN{
				splitAmount,
				splitOldTicker,
				splitOldCount,
				splitNewTicker,
				splitNewCount,
				splitDate,
			})
			//fmt.Println(result[1:])
		}

	}
	return resultList, nil
}

func stockOrderSplitReEvaluation(splits []SecurityIN, orders []Order) {
	for _, split := range splits {
		for index, order := range orders {

			if order.Date.After(split.Datetime) {
				// we dont need to recalculate orders after split date
				//fmt.Println("Orders are older, then current split, breaking inner loop")
				break
			}
			if order.Ticker != split.OldTicker {
				//fmt.Println("Orders ticker is different from split, thu skip order")
				continue
			}
			//fmt.Printf("Order %v should be reconsidered due to split %v", order, split)
			currentOrder := &orders[index]
			if currentOrder.Ticker != split.NewTicker {
				currentOrder.Ticker = split.NewTicker
			}
			currentOrder.Quantity = currentOrder.Quantity * split.NewCount
			currentOrder.Price = currentOrder.Price / float64(split.NewCount)

		}
	}
}

// func SplitsStockValidate(stockSplits []SecurityIN)  error {
//     splitMapAmount := make(map[string]int)
//     splitMapTranslator := make(map[string]SecurityIN)
//     for _, split := range(stockSplits) {
//         splitMapAmount[split.OldTicker] = splitMapAmount[split.OldTicker] + split.Quantity*split.Quantity
//         splitMapTranslator[split.OldTicker] = split
//     }
//
//     for key, v := range(splitMapAmount) {
//         if v !=0 {
//             return  errors.New(fmt.Sprintf("Amount mismatch for stock %v split", key))
//         }
//     }
//     return nil
// }
