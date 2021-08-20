package trading

import (
	"fmt"
	// "math"
	"github.com/neshcheretal/ffparser/jsonreport"
	"os"
	"time"
)

type Trade struct {
	OpenDate       time.Time
	CloseDate      time.Time
	OpenPrice      float64
	OpenUahRate    float64
	ClosePrice     float64
	CloseUahRate   float64
	Quantity       int
	OpenComission  float64
	CloseComission float64
}

func AllTradeListPrepare(t StockOrders) []Trade {
	madeTrades := make([]Trade, 0)
	buyQueue := t.Bought
	sellQueue := t.Sold
	for {
		if len(sellQueue) > 0 {
			var trade_quantity int
			var trade_buy_comission float64
			var trade_sell_comission float64
			sellTransaction := sellQueue[0]
			buyTransaction := buyQueue[0]
			if sellTransaction.Quantity > buyTransaction.Quantity {
				trade_quantity = buyTransaction.Quantity
				trade_buy_comission = buyTransaction.Comission
				trade_sell_comission = sellTransaction.Comission * float64(buyTransaction.Quantity) / float64(sellTransaction.Quantity)
				sellQueue[0] = Order{
					sellTransaction.Date,
					sellTransaction.Ticker,
					sellTransaction.Transaction,
					sellTransaction.Quantity - buyTransaction.Quantity,
					sellTransaction.Price,
					sellTransaction.UahRate,
					sellTransaction.Currency,
					sellTransaction.Comission - sellTransaction.Comission*float64(buyTransaction.Quantity)/float64(sellTransaction.Quantity),
				}
				buyQueue = buyQueue[1:]

			} else if sellTransaction.Quantity == buyTransaction.Quantity {
				trade_quantity = buyTransaction.Quantity
				trade_buy_comission = buyTransaction.Comission
				trade_sell_comission = sellTransaction.Comission
				// pop element from sell orders queue
				if len(sellQueue) > 1 {
					sellQueue = sellQueue[1:]
				} else {
					sellQueue = make([]Order, 0)
				}

				// pop element from sell orders queue
				if len(buyQueue) > 1 {
					buyQueue = buyQueue[1:]
				} else {
					buyQueue = make([]Order, 0)
				}

			} else if sellTransaction.Quantity < buyTransaction.Quantity {
				trade_quantity = sellTransaction.Quantity
				trade_buy_comission = buyTransaction.Comission * float64(sellTransaction.Quantity) / float64(buyTransaction.Quantity)
				trade_sell_comission = sellTransaction.Comission
				sellQueue = sellQueue[1:]
				buyQueue[0] = Order{
					buyTransaction.Date,
					buyTransaction.Ticker,
					buyTransaction.Transaction,
					buyTransaction.Quantity - sellTransaction.Quantity,
					buyTransaction.Price,
					buyTransaction.UahRate,
					buyTransaction.Currency,
					buyTransaction.Comission - buyTransaction.Comission*float64(sellTransaction.Quantity)/float64(buyTransaction.Quantity),
				}
			}
			madeTrade := Trade{
				buyTransaction.Date,
				sellTransaction.Date,
				buyTransaction.Price,
				buyTransaction.UahRate,
				sellTransaction.Price,
				sellTransaction.UahRate,
				trade_quantity,
				trade_buy_comission,
				trade_sell_comission,
			}
			madeTrades = append(madeTrades, madeTrade)
		} else {
			break
		}
	}
	return madeTrades
}

func FilterDateTrades(trades []Trade, startDate time.Time, endDate time.Time) []Trade {
	filteredTrades := make([]Trade, 0)
	if startDate.IsZero() && endDate.IsZero() {
		filteredTrades = trades
	} else if !startDate.IsZero() && endDate.IsZero() {
		for _, trade := range trades {
			if trade.CloseDate.After(startDate) {
				filteredTrades = append(filteredTrades, trade)
			}
		}
	} else if startDate.IsZero() && !endDate.IsZero() {
		for _, trade := range trades {
			if trade.CloseDate.Before(endDate) {
				filteredTrades = append(filteredTrades, trade)
			}
		}
	} else {
		for _, trade := range trades {
			if trade.CloseDate.After(startDate) && trade.CloseDate.Before(endDate) {
				filteredTrades = append(filteredTrades, trade)
			}
		}
	}
	return filteredTrades
}

//func StockTradePreparationWrapper(jsonResport jsonreport.JsonReport, startDate time.Time, endDate time.Time, ch  chan map[string][]trading.Trade) (map[string][]Trade, error) {
func StockTradePreparationWrapper(jsonResport jsonreport.JsonReport, startDate time.Time, endDate time.Time, ch chan map[string][]Trade) {

	stockOrderList, err := OrderListPreparator(jsonResport.Trades.Detailed)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stockSplits, err := SplitsStockParser(jsonResport.Securities_in_outs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	stockOrderSplitReEvaluation(stockSplits, stockOrderList)

	stockOrdersMap := OrdersStockParser(stockOrderList)

	stockTradesMap := make(map[string][]Trade)
	for ticker, stockTransactions := range stockOrdersMap {
		if len(stockTransactions.Sold) > 0 && stockTransactions.SoldCount > stockTransactions.BoughtCount {
			fmt.Print(ticker, " ticker was renamed(or report is incomplete), has to be handled manually \n")
		} else if len(stockTransactions.Sold) > 0 {
			allMadeTrades := AllTradeListPrepare(stockTransactions)
			dateFilteredTrades := FilterDateTrades(allMadeTrades, startDate, endDate)
			if len(dateFilteredTrades) != 0 {
				stockTradesMap[ticker] = dateFilteredTrades
			}
		}
	}
	ch <- stockTradesMap
}

// func tradeProfitCalculator(trademap map[string][]Trade) map[string]TradeProfit {
// 	stockProfit := make(map[string]TradeProfit)
// 	for ticker, stockTrades := range trademap {
// 		fmt.Print(ticker, ":\n")
// 		var ticker_profit float64
// 		for i, trade := range stockTrades {
// 			trade_profit := math.Floor((trade.ClosePrice-trade.OpenPrice)*float64(trade.Quantity)*100) / 100
// 			trade_comissions := math.Floor((trade.OpenComission+trade.CloseComission)*100) / 100
// 			tradeFinalProfit := math.Floor((trade_profit-trade_comissions)*100) / 100
// 			fmt.Printf("Trade %v: Profit is  %v\n", i, trade_profit)
// 			fmt.Printf("Trade %v: Total comission is %v\n", i, trade_comissions)
// 			fmt.Printf("Trade %v: Final Profit is `Profit - Comission` %v\n", i, tradeFinalProfit)
// 			ticker_profit += tradeFinalProfit
// 		}
// 		fmt.Print("Stock profit in USD is ", ticker_profit, "\n")
// 		stockProfit[ticker] = TradeProfit{
// 			ticker_profit,
// 		}
// 	}
// 	return stockProfit
// }
