package trading

import (
	"github.com/neshcheretal/ffparser/jsonreport"
	"github.com/neshcheretal/ffparser/nbu"
	"strings"
	"time"
)

// Type for broker transaction (separate stock buy or sell orders)
type Order struct {
	Date        time.Time
	PayDate     time.Time
	Ticker      string
	Transaction string
	Quantity    int
	Price       float64
	UahRate     float64
	Currency    string
	Comission   float64
}

// Custom type to group all orders for specific stock based on transaction type: buy/sell
type StockOrders struct {
	Bought      []Order
	BoughtCount int
	Sold        []Order
	SoldCount   int
}

func (s *StockOrders) setBought(newBought []Order) {
	s.Bought = newBought
}

func (s *StockOrders) setBoughtCount(newBoughtCount int) {
	s.BoughtCount = newBoughtCount
}

func (s *StockOrders) setSold(newSold []Order) {
	s.Sold = newSold
}

func (s *StockOrders) setSoldCount(newSoldCount int) {
	s.SoldCount = newSoldCount
}

//Receive pased execl report
func OrderListPreparator(jsonOrders []jsonreport.JsonOrder) ([]Order, error) {
	orders := make([]Order, 0)
	dateLayout := "2006-01-02 15:04:05"
	payDateLayout := "2006-01-02"
	for _, trade := range jsonOrders {
		orderDate, err := time.Parse(dateLayout, trade.Date)
		if err != nil {
			return nil, err
		}

		orderPayDate, err := time.Parse(payDateLayout, trade.Pay_d)
		if err != nil {
			return nil, err
		}

		ticker := trade.Instr_nm
		// Handle IPO tickers which initially created with .BLOCKED suffix that being dropped after lock-up period expiration
		if strings.HasSuffix(ticker, ".BLOCKED") {
			ticker = strings.ReplaceAll(ticker, ".BLOCKED", "")
		}
		// Handle stock suffix used by FF
		if strings.HasSuffix(ticker, ".US") {
			ticker = strings.ReplaceAll(ticker, ".US", "")
		}
		orderTicker := ticker
		orderTransaction := trade.Operation
		orderQuantity := trade.Q
		orderPrice := trade.P
		orderCurrency := trade.Curr_c
		orderComission := trade.Commission
		orderUsdPriceUah, err := nbu.GetConversionRates(orderPayDate, orderCurrency)
		if err != nil {
			return nil, err
		}

		newtransact := Order{
			orderDate,
			orderPayDate,
			orderTicker,
			orderTransaction,
			orderQuantity,
			orderPrice,
			orderUsdPriceUah,
			orderCurrency,
			orderComission,
		}
		orders = append(orders, newtransact)
	}
	return orders, nil
}

//Receive pased execl report
func OrdersStockParser(orders []Order) map[string]StockOrders {
	ordersMap := make(map[string]StockOrders)
	for _, order := range orders {
		currentTikerOrders := ordersMap[order.Ticker]
		if order.Transaction == "buy" {
			currentTikerBuyOrders := currentTikerOrders.Bought
			currentTikerBuyOrders = append(currentTikerBuyOrders, order)
			currentTikerOrders.setBought(currentTikerBuyOrders)
			currentTikerBuyOrdersCount := currentTikerOrders.BoughtCount
			currentTikerOrders.setBoughtCount(currentTikerBuyOrdersCount + order.Quantity)
		} else if order.Transaction == "sell" {
			currentTikerSellOrders := currentTikerOrders.Sold
			currentTikerSellOrders = append(currentTikerSellOrders, order)
			currentTikerOrders.setSold(currentTikerSellOrders)
			currentTikerSellOrdersCount := currentTikerOrders.SoldCount
			currentTikerOrders.setSoldCount(currentTikerSellOrdersCount + order.Quantity)
		}

		ordersMap[order.Ticker] = currentTikerOrders
	}
	return ordersMap
}
