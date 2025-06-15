package cmd

import (
	"fmt"
	"strings"
	"trader/internal/common"
	"trader/internal/config"
	"trader/internal/scraping"
	"trader/internal/service"
	"trader/internal/tools"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var securityCmd = &cobra.Command{
	Use:   "security",
	Short: "Tool to get security information",
	Long:  `Tool to get security information`,
}

var securityPurchaseBalanceByTickersCmd = &cobra.Command{
	Use:   "purchase-balance --stocks [tickers ...] --reits [tickers ...] --amount <float>",
	Short: "Purchase balance by tickers",
	Long:  `Purchase balance by tickers`,
	Run: func(cmd *cobra.Command, args []string) {
		config := config.NewConfig()
		stockScraping := scraping.NewStockScraping(config)
		stockService := service.NewStockService(stockScraping)
		reitScraping := scraping.NewReitScraping(config)
		reitService := service.NewReitService(reitScraping)
		purchaseBalanceService := service.NewPurchaseBalanceService(stockService, reitService)

		stocks := strings.Split(flagStocks, ",")
		reits := strings.Split(flagReits, ",")
		purchaseBalance := purchaseBalanceService.PurchaseBalancesBySecurities(stocks, reits, flagAmount)
		if len(purchaseBalance.SecuritiesBalance) == 0 {
			fmt.Println("Error: tickers not found!")
			return
		}
		t := common.NewTableWriter(flagNoColor)
		t.AppendHeader(table.Row{"TICKER", "TYPE", "PRICE", "COUNT", "TOTAL", "CURRENCY", "CAPTURED AT"})
		currency := purchaseBalance.SecuritiesBalance[0].Security.Currency.String()
		for _, purchase := range purchaseBalance.SecuritiesBalance {
			t.AppendRow(table.Row{purchase.Security.Ticker, purchase.Security.Type, tools.TableRowValue(purchase.Security.Price), purchase.Count, tools.TableRowValue(purchase.TotalAmount()), currency, tools.TableRowValue(purchase.Security.CapturedAt)})
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
		t.AppendFooter(table.Row{"", "", "", purchaseBalance.TotalCount(), tools.TableRowValue(purchaseBalance.AmountSpent()), currency, "SPENT AMOUNT"})
		t.AppendFooter(table.Row{"", "", "", "", tools.TableRowValue(purchaseBalance.RemainingBalance()), currency, "REMAINING AMOUNT"})
		t.SetIndexColumn(1)
		if flagCsv {
			t.RenderCSV()
		} else {
			t.Render()
		}
	},
}

func init() {
	rootCmd.AddCommand(securityCmd)
	securityCmd.AddCommand(securityPurchaseBalanceByTickersCmd)

	securityPurchaseBalanceByTickersCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	securityPurchaseBalanceByTickersCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")
	securityPurchaseBalanceByTickersCmd.Flags().StringVarP(&flagStocks, "stocks", "s", "", "List of stocks to purchase [ticker1,ticker2...] (required)")
	securityPurchaseBalanceByTickersCmd.Flags().StringVarP(&flagReits, "reits", "r", "", "List of REITs to purchase [ticker1,ticker2...] (required)")
	securityPurchaseBalanceByTickersCmd.Flags().Float64VarP(&flagAmount, "amount", "a", 0.0, "Amount invested (required)")
	securityPurchaseBalanceByTickersCmd.MarkFlagsRequiredTogether("amount")
}
