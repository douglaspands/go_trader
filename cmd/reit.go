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

var reitCmd = &cobra.Command{
	Use:   "reit",
	Short: "Tool to get reit information",
	Long:  `Tool to get reit information`,
}

var reitGetReitByTickerCmd = &cobra.Command{
	Use:   "get [ticker]",
	Short: "Get a reit by ticker",
	Long:  `Get a reit by ticker`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticker := args[0]
		reit := service.GetReit(ticker)
		if reit == nil {
			fmt.Printf("Error: ticker \"%s\" not found!\n", ticker)
			return
		}
		t := common.NewTableWriter(flagNoColor)
		t.AppendHeader(table.Row{"FIELD", "VALUE"})
		t.AppendRow(table.Row{"Ticker", reit.Ticker})
		t.AppendRow(table.Row{"Name", reit.Name})
		t.AppendRow(table.Row{"Admin", reit.Admin})
		t.AppendRow(table.Row{"Document", reit.Document})
		t.AppendRow(table.Row{"Segment", reit.Segment})
		t.AppendRow(table.Row{"Currency", reit.Currency.String()})
		t.AppendRow(table.Row{"Price", tools.TableRowValue(reit.Price)})
		t.AppendRow(table.Row{"CapturedAt", tools.TableRowValue(reit.CapturedAt)})
		t.AppendRow(table.Row{"Origin", reit.Origin})
		// t.AppendRow(table.Row{"Description", reit.Description})
		t.SetIndexColumn(1)
		if flagCsv {
			t.RenderCSV()
		} else {
			t.Render()
		}
	},
}

var reitListReitsByTickersCmd = &cobra.Command{
	Use:   "list [tickers ...]",
	Short: "List reits by tickers",
	Long:  `List reits by tickers`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tickers := args
		reits := service.ListReits(tickers)
		if len(reits) == 0 {
			fmt.Println("Error: tickers not found!")
			return
		}
		t := common.NewTableWriter(flagNoColor)
		t.AppendHeader(table.Row{"TICKER", "NAME", "DOCUMENT", "PRICE", "CURRENCY", "CAPTURED AT"})
		for _, reit := range reits {
			t.AppendRow(table.Row{reit.Ticker, reit.Name, reit.Document, tools.TableRowValue(reit.Price), reit.Currency.String(), tools.TableRowValue(reit.CapturedAt)})
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

var reitPurchaseBalanceByTickersCmd = &cobra.Command{
	Use:   "purchase-balance [tickers ...] --amount <float>",
	Short: "Purchase balance by tickers",
	Long:  `Purchase balance by tickers`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tickers := args
		purchaseBalance := service.MakeReitPurchaseBalance(tickers, flagAmount)
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
	rootCmd.AddCommand(reitCmd)
	reitCmd.AddCommand(reitGetReitByTickerCmd)
	reitCmd.AddCommand(reitListReitsByTickersCmd)
	reitCmd.AddCommand(reitPurchaseBalanceByTickersCmd)

	reitGetReitByTickerCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	reitGetReitByTickerCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")

	reitListReitsByTickersCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	reitListReitsByTickersCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")

	reitPurchaseBalanceByTickersCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	reitPurchaseBalanceByTickersCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")
	reitPurchaseBalanceByTickersCmd.Flags().Float64VarP(&flagAmount, "amount", "a", 0.0, "Amount invested (required)")
	reitPurchaseBalanceByTickersCmd.MarkFlagsRequiredTogether("amount")
}
