package cmd

import (
	"fmt"
	"trader/internal/common"
	"trader/internal/config"
	"trader/internal/scraping"
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
		config := config.NewConfig()
		reitScraping := scraping.NewReitScraping(config)
		reitService := service.NewReitService(reitScraping)

		ticker := args[0]
		reit := reitService.GetReitByTicker(ticker)
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
		config := config.NewConfig()
		reitScraping := scraping.NewReitScraping(config)
		reitService := service.NewReitService(reitScraping)

		tickers := args
		reits := reitService.ListReitsByTickers(tickers)
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

func init() {
	rootCmd.AddCommand(reitCmd)
	reitCmd.AddCommand(reitGetReitByTickerCmd)
	reitCmd.AddCommand(reitListReitsByTickersCmd)

	reitGetReitByTickerCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	reitGetReitByTickerCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")

	reitListReitsByTickersCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	reitListReitsByTickersCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")
}
