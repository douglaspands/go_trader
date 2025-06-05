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
		t := common.NewTableWriter()
		t.AppendRows(tools.StructToTableRowsFieldValue(reit, []string{"Description"}))
		t.SetIndexColumn(1)
		t.Render()
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
		t := common.NewTableWriter()
		for idx, reit := range reits {
			if idx == 0 {
				t.AppendHeader(tools.StructToTableHeader(reit, []string{"Description", "Origin"}))
			}
			t.AppendRow(tools.StructToTableRow(reit, []string{"Description", "Origin"}))
		}
		t.SetColumnConfigs([]table.ColumnConfig{
			{
				Name:        "PRICE",
				AlignHeader: text.AlignRight,
				Align:       text.AlignRight,
			},
		})
		t.SetIndexColumn(1)
		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(reitCmd)
	reitCmd.AddCommand(reitGetReitByTickerCmd)
	reitCmd.AddCommand(reitListReitsByTickersCmd)
}
