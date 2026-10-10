package cmd

import (
	"errors"
	"trader/internal/common"
	"trader/internal/service"
	"trader/internal/tools"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

type SecurityCommand interface {
	InitApp(rootCmd RootCommand)
}

type securityCommand struct {
	purchaseBalanceService service.PurchaseBalanceService
	// Commands
	rootCmd *cobra.Command
	// Flags
	flagAmount float64
	noColor    bool
	csv        bool
	flagStocks string
	flagReits  string
}

func (sc *securityCommand) purchaseBalanceByTickersCmd(cmd *cobra.Command, args []string) error {
	if err := validateAmount(sc.flagAmount); err != nil {
		return err
	}
	stocks := tools.SplitList(sc.flagStocks, ",")
	reits := tools.SplitList(sc.flagReits, ",")
	purchaseBalance, err := sc.purchaseBalanceService.PurchaseBalancesBySecurities(stocks, reits, sc.flagAmount)
	if err != nil {
		cmd.SilenceUsage = true
		return cleanError(err)
	}
	if len(purchaseBalance.SecuritiesBalance) == 0 {
		cmd.SilenceUsage = true
		return errors.New("tickers not found!")
	}
	t := common.NewTableWriter(sc.noColor, cmd.OutOrStdout())
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
	render(t, sc.csv)
	return nil
}

func (sc *securityCommand) setup() {
	purchaseBalanceByTickersCmd := &cobra.Command{
		Use:   "purchase-balance --stocks [tickers ...] --reits [tickers ...] --amount <float>",
		Short: "Purchase balance by tickers",
		Long:  `Purchase balance by tickers`,
		RunE:  sc.purchaseBalanceByTickersCmd,
	}
	purchaseBalanceByTickersCmd.Flags().BoolVar(&sc.noColor, "no-color", false, "Output without color")
	purchaseBalanceByTickersCmd.Flags().BoolVar(&sc.csv, "csv", false, "Output csv format")
	purchaseBalanceByTickersCmd.Flags().StringVarP(&sc.flagStocks, "stocks", "s", "", "List of stocks to purchase [ticker1,ticker2...] (required)")
	purchaseBalanceByTickersCmd.Flags().StringVarP(&sc.flagReits, "reits", "r", "", "List of REITs to purchase [ticker1,ticker2...] (required)")
	purchaseBalanceByTickersCmd.Flags().Float64VarP(&sc.flagAmount, "amount", "a", 0.0, "Amount invested (required)")
	purchaseBalanceByTickersCmd.MarkFlagRequired("amount")
	sc.rootCmd.AddCommand(purchaseBalanceByTickersCmd)
}

func (sc *securityCommand) InitApp(rootCmd RootCommand) {
	rootCmd.GetCobraCommand().AddCommand(sc.rootCmd)
	sc.setup()
}

func NewSecurityCommand(purchaseBalanceService service.PurchaseBalanceService) SecurityCommand {
	return &securityCommand{
		purchaseBalanceService: purchaseBalanceService,
		rootCmd: &cobra.Command{
			Use:   "security",
			Short: "Tool to get security information",
			Long:  `Tool to get security information`,
		},
	}
}
