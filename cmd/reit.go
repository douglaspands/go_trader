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

type ReitCommand interface {
	InitApp(rootCmd RootCommand)
}

type reitCommand struct {
	reitService service.ReitService
	// Commands
	rootCmd               *cobra.Command
	getReitByTickerCmd    *cobra.Command
	listReitsByTickersCmd *cobra.Command
}

func (rc *reitCommand) setup() {

	// getReitByTickerCmd
	rc.getReitByTickerCmd = &cobra.Command{
		Use:   "get [ticker]",
		Short: "Get a reit by ticker",
		Long:  `Get a reit by ticker`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ticker := args[0]
			reit := rc.reitService.GetReitByTicker(ticker)
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
	rc.getReitByTickerCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	rc.getReitByTickerCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")

	// listReitsByTickersCmd
	rc.listReitsByTickersCmd = &cobra.Command{
		Use:   "list [tickers ...]",
		Short: "List reits by tickers",
		Long:  `List reits by tickers`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tickers := args
			reits := rc.reitService.ListReitsByTickers(tickers)
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
	rc.listReitsByTickersCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	rc.listReitsByTickersCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")

}

func (rc *reitCommand) register() {
	rc.rootCmd.AddCommand(rc.getReitByTickerCmd)
	rc.rootCmd.AddCommand(rc.listReitsByTickersCmd)
}

func (rc *reitCommand) InitApp(rootCmd RootCommand) {
	rootCmd.GetCobraCommand().AddCommand(rc.rootCmd)
	rc.setup()
	rc.register()
}

func NewReitCommand(reitService service.ReitService) ReitCommand {
	return &reitCommand{
		reitService: reitService,
		rootCmd: &cobra.Command{
			Use:   "reit",
			Short: "Tool to get reit information",
			Long:  `Tool to get reit information`,
		},
	}
}
