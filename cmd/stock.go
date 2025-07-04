package cmd

import (
	"fmt"
	"trader/internal/common"
	"trader/internal/service"
	"trader/internal/tools"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var stockCmd = &cobra.Command{
	Use:   "stock",
	Short: "Tool to get stock information",
	Long:  `Tool to get stock information`,
}

var stockGetStockByTickerCmd = &cobra.Command{
	Use:   "get [ticker]",
	Short: "Get a stock by ticker",
	Long:  `Get a stock by ticker`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticker := args[0]
		stock := service.GetStock(ticker)
		if stock == nil {
			fmt.Printf("Error: ticker \"%s\" not found!\n", ticker)
			return
		}
		t := common.NewTableWriter(flagNoColor)
		t.AppendHeader(table.Row{"FIELD", "VALUE"})
		t.AppendRow(table.Row{"Ticker", stock.Ticker})
		t.AppendRow(table.Row{"Name", stock.Name})
		t.AppendRow(table.Row{"Document", stock.Document})
		t.AppendRow(table.Row{"Currency", stock.Currency.String()})
		t.AppendRow(table.Row{"Price", tools.TableRowValue(stock.Price)})
		t.AppendRow(table.Row{"CapturedAt", tools.TableRowValue(stock.CapturedAt)})
		t.AppendRow(table.Row{"Origin", stock.Origin})
		// t.AppendRow(table.Row{"Description", stock.Description})
		t.SetIndexColumn(1)
		if flagCsv {
			t.RenderCSV()
		} else {
			t.Render()
		}
	},
}

var stockListStocksByTickersCmd = &cobra.Command{
	Use:   "list [tickers ...]",
	Short: "List stocks by tickers",
	Long:  `List stocks by tickers`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tickers := args
		stocks := service.ListStocks(tickers)
		if len(stocks) == 0 {
			fmt.Println("Error: tickers not found!")
			return
		}
		t := common.NewTableWriter(flagNoColor)
		t.AppendHeader(table.Row{"TICKER", "NAME", "DOCUMENT", "PRICE", "CURRENCY", "CAPTURED AT"})
		for _, stock := range stocks {
			t.AppendRow(table.Row{stock.Ticker, stock.Name, stock.Document, tools.TableRowValue(stock.Price), stock.Currency.String(), tools.TableRowValue(stock.CapturedAt)})
		}
		t.SetColumnConfigs([]table.ColumnConfig{
			{
				Name:        "TICKER",
				AlignHeader: text.AlignRight,
				Align:       text.AlignRight,
			},
			{
				Name:        "PRICE",
				AlignHeader: text.AlignRight,
				Align:       text.AlignRight,
			},
		})
		t.SetIndexColumn(1)
		if flagCsv {
			t.RenderCSV()
		} else {
			t.Render()
		}
	},
}

var stockPurchaseBalanceByTickersCmd = &cobra.Command{
	Use:   "purchase-balance [tickers ...] --amount <float>",
	Short: "Purchase balance by tickers",
	Long:  `Purchase balance by tickers`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tickers := args
		purchaseBalance := service.MakeStockPurchaseBalance(tickers, flagAmount)
		if len(purchaseBalance.SecuritiesBalance) == 0 {
			fmt.Println("Error: tickers not found!")
			return
		}
		t := common.NewTableWriter(flagNoColor)
		t.AppendHeader(table.Row{"TICKER", "PRICE", "COUNT", "TOTAL", "CURRENCY", "CAPTURED AT"})
		currency := purchaseBalance.SecuritiesBalance[0].Security.Currency.String()
		for _, purchase := range purchaseBalance.SecuritiesBalance {
			t.AppendRow(table.Row{purchase.Security.Ticker, tools.TableRowValue(purchase.Security.Price), purchase.Count, tools.TableRowValue(purchase.TotalAmount()), currency, tools.TableRowValue(purchase.Security.CapturedAt)})
		}
		t.SetColumnConfigs([]table.ColumnConfig{
			{
				Name:        "TICKER",
				Align:       text.AlignRight,
				AlignHeader: text.AlignRight,
				AlignFooter: text.AlignRight,
			},
			{
				Name:        "PRICE",
				Align:       text.AlignRight,
				AlignHeader: text.AlignRight,
				AlignFooter: text.AlignRight,
			},
			{
				Name:        "TOTAL",
				Align:       text.AlignRight,
				AlignHeader: text.AlignRight,
				AlignFooter: text.AlignRight,
			},
		})
		t.AppendFooter(table.Row{"", "", purchaseBalance.TotalCount(), tools.TableRowValue(purchaseBalance.AmountSpent()), currency, "SPENT AMOUNT"})
		t.AppendFooter(table.Row{"", "", "", tools.TableRowValue(purchaseBalance.RemainingBalance()), currency, "REMAINING AMOUNT"})
		t.SetIndexColumn(1)
		if flagCsv {
			t.RenderCSV()
		} else {
			t.Render()
		}
	},
}

func init() {
	rootCmd.AddCommand(stockCmd)
	stockCmd.AddCommand(stockGetStockByTickerCmd)
	stockCmd.AddCommand(stockListStocksByTickersCmd)
	stockCmd.AddCommand(stockPurchaseBalanceByTickersCmd)

	stockGetStockByTickerCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	stockGetStockByTickerCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")

	stockListStocksByTickersCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	stockListStocksByTickersCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")

	stockPurchaseBalanceByTickersCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	stockPurchaseBalanceByTickersCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")
	stockPurchaseBalanceByTickersCmd.Flags().Float64VarP(&flagAmount, "amount", "a", 0.0, "Amount invested (required)")
	stockPurchaseBalanceByTickersCmd.MarkFlagsRequiredTogether("amount")
}
