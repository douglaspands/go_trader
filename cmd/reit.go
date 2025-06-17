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
	reitService            service.ReitService
	purchaseBalanceService service.PurchaseBalanceService
	// Commands
	rootCmd *cobra.Command
	// Flags
	flagAmount float64
}

func (rc *reitCommand) getReitByTickerCmd(cmd *cobra.Command, args []string) {
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
}

func (rc *reitCommand) listReitsByTickersCmd(cmd *cobra.Command, args []string) {
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
}

func (rc *reitCommand) purchaseBalanceByTickersCmd(cmd *cobra.Command, args []string) {
	tickers := args
	purchaseBalance := rc.purchaseBalanceService.PurchaseBalancesBySecurities([]string{}, tickers, rc.flagAmount)
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
}

func (rc *reitCommand) setup() {
	getReitByTickerCmd := &cobra.Command{
		Use:   "get [ticker]",
		Short: "Get a reit by ticker",
		Long:  `Get a reit by ticker`,
		Args:  cobra.ExactArgs(1),
		Run:   rc.getReitByTickerCmd,
	}
	getReitByTickerCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	getReitByTickerCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")
	rc.rootCmd.AddCommand(getReitByTickerCmd)

	listReitsByTickersCmd := &cobra.Command{
		Use:   "list [tickers ...]",
		Short: "List reits by tickers",
		Long:  `List reits by tickers`,
		Args:  cobra.MinimumNArgs(1),
		Run:   rc.listReitsByTickersCmd,
	}
	listReitsByTickersCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	listReitsByTickersCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")
	rc.rootCmd.AddCommand(listReitsByTickersCmd)

	purchaseBalanceByTickersCmd := &cobra.Command{
		Use:   "purchase-balance [tickers ...] --amount <float>",
		Short: "Purchase balance by tickers",
		Long:  `Purchase balance by tickers`,
		Args:  cobra.MinimumNArgs(1),
		Run:   rc.purchaseBalanceByTickersCmd,
	}
	purchaseBalanceByTickersCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Output without color")
	purchaseBalanceByTickersCmd.Flags().BoolVar(&flagCsv, "csv", false, "Output csv format")
	purchaseBalanceByTickersCmd.Flags().Float64VarP(&rc.flagAmount, "amount", "a", 0.0, "Amount invested (required)")
	purchaseBalanceByTickersCmd.MarkFlagsRequiredTogether("amount")
	rc.rootCmd.AddCommand(purchaseBalanceByTickersCmd)
}

func (rc *reitCommand) InitApp(rootCmd RootCommand) {
	rootCmd.GetCobraCommand().AddCommand(rc.rootCmd)
	rc.setup()
}

func NewReitCommand(reitService service.ReitService, purchaseBalanceService service.PurchaseBalanceService) ReitCommand {
	return &reitCommand{
		reitService:            reitService,
		purchaseBalanceService: purchaseBalanceService,
		rootCmd: &cobra.Command{
			Use:   "reit",
			Short: "Tool to get reit information",
			Long:  `Tool to get reit information`,
		},
	}
}
