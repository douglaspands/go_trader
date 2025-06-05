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
		t := common.NewTableWriter()
		t.AppendRows(tools.StructToTableRowsFieldValue(stock, []string{"Description"}))
		t.Render()
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
		t := common.NewTableWriter()
		for idx, stock := range stocks {
			if idx == 0 {
				t.AppendHeader(tools.StructToTableHeader(stock, []string{"Description", "Origin"}))
			}
			t.AppendRow(tools.StructToTableRow(stock, []string{"Description", "Origin"}))
		}
		t.SetColumnConfigs([]table.ColumnConfig{
			{
				Name:        "PRICE",
				AlignHeader: text.AlignRight,
				Align:       text.AlignRight,
			},
		})
		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(stockCmd)
	stockCmd.AddCommand(stockGetStockByTickerCmd)
	stockCmd.AddCommand(stockListStocksByTickersCmd)
}
