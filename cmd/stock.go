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

type StockCommand interface {
	InitApp(rootCmd RootCommand)
}

type stockCommand struct {
	stockService           service.StockService
	purchaseBalanceService service.PurchaseBalanceService
	// Commands
	rootCmd *cobra.Command
	// Flags
	flagAmount float64
	noColor    bool
	csv        bool
}

func (sc *stockCommand) getStockByTickerCmd(cmd *cobra.Command, args []string) error {
	ticker := args[0]
	stock := sc.stockService.GetStockByTicker(ticker)
	if stock == nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("ticker \"%s\" not found!", ticker)
	}
	t := common.NewTableWriter(sc.noColor, cmd.OutOrStdout())
	t.AppendHeader(table.Row{"FIELD", "VALUE"})
	t.AppendRow(table.Row{"Ticker", stock.Ticker})
	t.AppendRow(table.Row{"Name", cell(stock.Name, sc.csv)})
	t.AppendRow(table.Row{"Document", cell(stock.Document, sc.csv)})
	t.AppendRow(table.Row{"Currency", stock.Currency.String()})
	t.AppendRow(table.Row{"Price", tools.TableRowValue(stock.Price)})
	t.AppendRow(table.Row{"CapturedAt", tools.TableRowValue(stock.CapturedAt)})
	t.AppendRow(table.Row{"Origin", stock.Origin})
	// t.AppendRow(table.Row{"Description", stock.Description})
	t.SetIndexColumn(1)
	render(t, sc.csv)
	return nil
}

func (sc *stockCommand) listStocksByTickersCmd(cmd *cobra.Command, args []string) error {
	tickers := args
	stocks, failures := sc.stockService.ListStocksByTickers(tickers)
	if len(stocks) == 0 {
		cmd.SilenceUsage = true
		return errors.New("tickers not found!")
	}
	t := common.NewTableWriter(sc.noColor, cmd.OutOrStdout())
	t.AppendHeader(table.Row{"TICKER", "NAME", "DOCUMENT", "PRICE", "CURRENCY", "CAPTURED AT"})
	for _, stock := range stocks {
		t.AppendRow(table.Row{stock.Ticker, cell(stock.Name, sc.csv), cell(stock.Document, sc.csv), tools.TableRowValue(stock.Price), stock.Currency.String(), tools.TableRowValue(stock.CapturedAt)})
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
	render(t, sc.csv)
	printWarnings(cmd.ErrOrStderr(), failures)
	return nil
}

func (sc *stockCommand) purchaseBalanceByTickersCmd(cmd *cobra.Command, args []string) error {
	if err := validateAmount(sc.flagAmount); err != nil {
		return err
	}
	tickers := args
	purchaseBalance, err := sc.purchaseBalanceService.PurchaseBalancesBySecurities(tickers, []string{}, sc.flagAmount)
	if err != nil {
		cmd.SilenceUsage = true
		return cleanError(err)
	}
	if len(purchaseBalance.SecuritiesBalance) == 0 {
		cmd.SilenceUsage = true
		return errors.New("tickers not found!")
	}
	t := common.NewTableWriter(sc.noColor, cmd.OutOrStdout())
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
	render(t, sc.csv)
	return nil
}

func (sc *stockCommand) setup() {
	getStockByTickerCmd := &cobra.Command{
		Use:   "get [ticker]",
		Short: "Get a stock by ticker",
		Long:  `Get a stock by ticker`,
		Args:  cobra.ExactArgs(1),
		RunE:  sc.getStockByTickerCmd,
	}
	getStockByTickerCmd.Flags().BoolVar(&sc.noColor, "no-color", false, "Output without color")
	getStockByTickerCmd.Flags().BoolVar(&sc.csv, "csv", false, "Output csv format")
	sc.rootCmd.AddCommand(getStockByTickerCmd)

	listStocksByTickersCmd := &cobra.Command{
		Use:   "list [tickers ...]",
		Short: "List stocks by tickers",
		Long:  `List stocks by tickers`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  sc.listStocksByTickersCmd,
	}
	listStocksByTickersCmd.Flags().BoolVar(&sc.noColor, "no-color", false, "Output without color")
	listStocksByTickersCmd.Flags().BoolVar(&sc.csv, "csv", false, "Output csv format")
	sc.rootCmd.AddCommand(listStocksByTickersCmd)

	purchaseBalanceByTickersCmd := &cobra.Command{
		Use:   "purchase-balance [tickers ...] --amount <float>",
		Short: "Purchase balance by tickers",
		Long:  `Purchase balance by tickers`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  sc.purchaseBalanceByTickersCmd,
	}
	purchaseBalanceByTickersCmd.Flags().BoolVar(&sc.noColor, "no-color", false, "Output without color")
	purchaseBalanceByTickersCmd.Flags().BoolVar(&sc.csv, "csv", false, "Output csv format")
	purchaseBalanceByTickersCmd.Flags().Float64VarP(&sc.flagAmount, "amount", "a", 0.0, "Amount invested (required)")
	purchaseBalanceByTickersCmd.MarkFlagRequired("amount")
	sc.rootCmd.AddCommand(purchaseBalanceByTickersCmd)
}

func (sc *stockCommand) InitApp(rootCmd RootCommand) {
	rootCmd.GetCobraCommand().AddCommand(sc.rootCmd)
	sc.setup()
}

func NewStockCommand(stockService service.StockService, purchaseBalanceService service.PurchaseBalanceService) StockCommand {
	return &stockCommand{
		stockService:           stockService,
		purchaseBalanceService: purchaseBalanceService,
		rootCmd: &cobra.Command{
			Use:   "stock",
			Short: "Tool to get stock information",
			Long:  `Tool to get stock information`,
		},
	}
}
