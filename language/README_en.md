### FF Broker Report Parser
#### About the program

This program is designed to process brokerage reports from Freedom Finance Europe in order to prepare tax reports for the State Tax Service of Ukraine. In accordance with article 170, paragraph 2 of the Tax Code of Ukraine: "Investment profit is calculated as the positive difference between the income received by the taxpayer from the sale of a separate investment asset, taking into account the exchange rate difference (if any), and its value, which is determined from the amount of documented acquisition costs such an asset, taking into account the norms of sub-clauses 170.2.4-170.2.6 of this clause (except for operations with derivatives) ".
Despite the fact that brokerage reports provide detailed information about the investment operations made, they do not take into account the currency exchange rate difference at the date of the purchase/sale of the asset. Also when displaying closed deals, the broker indicates the average price of purchases of the asset for the entire time, which does not allow to correctly determine the price difference and does not take into account the amount of brokerage commission, while according to the above requirement it must be taken into account in the calculation. This program has built-in integration with the API of the National Bank of Ukraine to obtain current exchange rates by the date of purchase/sale and payment of dividends.

#### Usage and launch:
The command is an executable binary file and is intended to be run from the command line. The launch options are listed below:
   - -end End day for calculating the made deals in the YYYY-MM-DD format.
   - -start Start day for calculating the made deals in the YYYY-MM-DD format.
   - -lang Report language. Supported languages are: UA, RU, EN. (default "EN")
   - -output Name of the xlsx file with results ("tax_calculation.xlsx" by default)
   - -report Path to the JSON file of the FF broker report. (Required)

Launch exaple for MacOS/Linux
> ./ffparser -report=full-report.json -start=2020-08-01 -lang=UA -output=results.xlsx

#### How to get brokerage report
In order to get brokerage report you have to log in to Freedom Finance Europe trading terminal and the follow the next sequence:

Menu > BROKERAGE REPORTS > BROKER REPORT > Report for the period

The period MUST be set from the day when brokerage account was open. Then chose a JSON format and download the file.

#### Calculation features
For the correct calculation of the results of transactions, regardless of the period that will be declared, it is necessary to use a full brokerage report for the entire period of work starting from the day the brokerage account was opened, for example, an asset sold in 2021 could be purchased in 2020 and all data is required to calculate a closed transaction. The program calculates the income from the saling of stocks and recived dividends. The tax amount is calculated acorfing to following rates:
- 18% Personal income tax
- 1.5% Military tax
- 9% Dividend income

It is important to note that the program works correctly only with "long" positions if there is a "short" positions made, the calculation may be incorrect. Also it should be noted that it does not take into account the declaration requirements given in subparagraphs [170.2.4-170.2.6] (https://zakon.rada.gov.ua/laws/show/2755-17?lang=ru#n3998) and in the case of such transactions, the final sum to declare must be reviewed separately.

#### Calculation format
The calculation results are provided as an xlsx file with three tabs:
- Trading - the first tab with the results of income calculation from stock trading
- Dividends - the second tab with a list of dividends accrued for the period
- Taxes - the third tab with the calculation of tax liabilities

It should be mentioned that if there is a stock split ocures inside calculated period, all operations for the stock are recalculated and displayed with values ​​after the split, if the stock ticker has been renamed, all operations on it will be displayed with the value of the ticker after the split. For example:
- On January 1, 2 TEST stocks were bought for $ 300
- On January 2, TEST stock worth $ 300 was sold
- On January 3, There was a split 1 TEST 300 $ -> 3 BEST 100 $
- On January 4, 1 TEST stock worth $ 100 was sold

The final report will display two closed trades
| Ticker | Purchase date | Quantity | Purchase price | Sale date | Selling price |
--------|--------------|------------|--------------|--------------|--------------|
| BEST | 1 January | 3 | 100$ | 2 January | 100$ |
| BEST | 1 January | 1 | 100$ | 2 January | 100$ |
