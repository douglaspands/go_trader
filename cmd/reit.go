package cmd

import (
	"errors"
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
	noColor    bool
	csv        bool
}

func (rc *reitCommand) getReitByTickerCmd(cmd *cobra.Command, args []string) error {
	ticker := args[0]
	reit := rc.reitService.GetReitByTicker(ticker)
	if reit == nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("ticker \"%s\" not found!", ticker)
	}
	t := common.NewTableWriter(rc.noColor, cmd.OutOrStdout())
	t.AppendHeader(table.Row{"FIELD", "VALUE"})
	t.AppendRow(table.Row{"Ticker", reit.Ticker})
	t.AppendRow(table.Row{"Name", cell(reit.Name, rc.csv)})
	t.AppendRow(table.Row{"Admin", cell(reit.Admin, rc.csv)})
	t.AppendRow(table.Row{"Document", cell(reit.Document, rc.csv)})
	t.AppendRow(table.Row{"Segment", cell(reit.Segment, rc.csv)})
	t.AppendRow(table.Row{"Currency", reit.Currency.String()})
	t.AppendRow(table.Row{"Price", tools.TableRowValue(reit.Price)})
	t.AppendRow(table.Row{"CapturedAt", tools.TableRowValue(reit.CapturedAt)})
	t.AppendRow(table.Row{"Origin", reit.Origin})
	// t.AppendRow(table.Row{"Description", reit.Description})
	t.SetIndexColumn(1)
	render(t, rc.csv)
	return nil
}

func (rc *reitCommand) listReitsByTickersCmd(cmd *cobra.Command, args []string) error {
	tickers := args
	reits, failures := rc.reitService.ListReitsByTickers(tickers)
	if len(reits) == 0 {
		cmd.SilenceUsage = true
		return errors.New("tickers not found!")
	}
	t := common.NewTableWriter(rc.noColor, cmd.OutOrStdout())
	t.AppendHeader(table.Row{"TICKER", "NAME", "DOCUMENT", "PRICE", "CURRENCY", "CAPTURED AT"})
	for _, reit := range reits {
		t.AppendRow(table.Row{reit.Ticker, cell(reit.Name, rc.csv), cell(reit.Document, rc.csv), tools.TableRowValue(reit.Price), reit.Currency.String(), tools.TableRowValue(reit.CapturedAt)})
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
	render(t, rc.csv)
	printWarnings(cmd.ErrOrStderr(), failures)
	return nil
}

func (rc *reitCommand) purchaseBalanceByTickersCmd(cmd *cobra.Command, args []string) error {
	if err := validateAmount(rc.flagAmount); err != nil {
		return err
	}
	tickers := args
	purchaseBalance, err := rc.purchaseBalanceService.PurchaseBalancesBySecurities([]string{}, tickers, rc.flagAmount)
	if err != nil {
		cmd.SilenceUsage = true
		return cleanError(err)
	}
	if len(purchaseBalance.SecuritiesBalance) == 0 {
		cmd.SilenceUsage = true
		return errors.New("tickers not found!")
	}
	t := common.NewTableWriter(rc.noColor, cmd.OutOrStdout())
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
	render(t, rc.csv)
	return nil
}

func (rc *reitCommand) setup() {
	getReitByTickerCmd := &cobra.Command{
		Use:   "get [ticker]",
		Short: "Get a reit by ticker",
		Long:  `Get a reit by ticker`,
		Args:  cobra.ExactArgs(1),
		RunE:  rc.getReitByTickerCmd,
	}
	getReitByTickerCmd.Flags().BoolVar(&rc.noColor, "no-color", false, "Output without color")
	getReitByTickerCmd.Flags().BoolVar(&rc.csv, "csv", false, "Output csv format")
	rc.rootCmd.AddCommand(getReitByTickerCmd)

	listReitsByTickersCmd := &cobra.Command{
		Use:   "list [tickers ...]",
		Short: "List reits by tickers",
		Long:  `List reits by tickers`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  rc.listReitsByTickersCmd,
	}
	listReitsByTickersCmd.Flags().BoolVar(&rc.noColor, "no-color", false, "Output without color")
	listReitsByTickersCmd.Flags().BoolVar(&rc.csv, "csv", false, "Output csv format")
	rc.rootCmd.AddCommand(listReitsByTickersCmd)

	purchaseBalanceByTickersCmd := &cobra.Command{
		Use:   "purchase-balance [tickers ...] --amount <float>",
		Short: "Purchase balance by tickers",
		Long:  `Purchase balance by tickers`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  rc.purchaseBalanceByTickersCmd,
	}
	purchaseBalanceByTickersCmd.Flags().BoolVar(&rc.noColor, "no-color", false, "Output without color")
	purchaseBalanceByTickersCmd.Flags().BoolVar(&rc.csv, "csv", false, "Output csv format")
	purchaseBalanceByTickersCmd.Flags().Float64VarP(&rc.flagAmount, "amount", "a", 0.0, "Amount invested (required)")
	purchaseBalanceByTickersCmd.MarkFlagRequired("amount")
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
