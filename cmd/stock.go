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

type StockCommand interface {
	InitApp(rootCmd RootCommand)
}

type stockCommand struct {
	stockService service.StockService
	// Commands
	rootCmd                *cobra.Command
	getStockByTickerCmd    *cobra.Command
	listStocksByTickersCmd *cobra.Command
}

func (sc *stockCommand) setup() {
	// getStockByTickerCmd
	sc.getStockByTickerCmd = &cobra.Command{
		Use:   "get [ticker]",
		Short: "Get a stock by ticker",
		Long:  `Get a stock by ticker`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ticker := args[0]
			stock := sc.stockService.GetStockByTicker(ticker)
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

	sc.getStockByTickerCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	sc.getStockByTickerCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")

	// stockListStocksByTickersCmd
	sc.listStocksByTickersCmd = &cobra.Command{
		Use:   "list [tickers ...]",
		Short: "List stocks by tickers",
		Long:  `List stocks by tickers`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tickers := args
			stocks := sc.stockService.ListStocksByTickers(tickers)
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

	sc.listStocksByTickersCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	sc.listStocksByTickersCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")

}

func (sc *stockCommand) register() {
	sc.rootCmd.AddCommand(sc.getStockByTickerCmd)
	sc.rootCmd.AddCommand(sc.listStocksByTickersCmd)
}

func (sc *stockCommand) InitApp(rootCmd RootCommand) {
	rootCmd.GetCobraCommand().AddCommand(sc.rootCmd)
	sc.setup()
	sc.register()
}

func NewStockCommand(stockService service.StockService) StockCommand {
	return &stockCommand{
		stockService: stockService,
		rootCmd: &cobra.Command{
			Use:   "stock",
			Short: "Tool to get stock information",
			Long:  `Tool to get stock information`,
		},
	}
}
