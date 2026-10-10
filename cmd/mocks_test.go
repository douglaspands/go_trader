package cmd_test

import (
	"bytes"
	"time"
	"trader/cmd"
	"trader/internal/resource"
)

var capturedAt = time.Date(2025, 6, 8, 22, 59, 37, 0, time.Local)

var brl = &resource.Currency{Code: "BRL", Sign: "R$", Description: "Brazilian Real"}

func newStock(ticker string, price float64) *resource.Security {
	return &resource.Security{
		Ticker:     ticker,
		Type:       resource.STOCK_TYPE,
		Name:       "Name " + ticker,
		Document:   "11.111.111/0001-11",
		Currency:   brl,
		Price:      price,
		Origin:     "https://example.com/acoes/" + ticker,
		CapturedAt: capturedAt,
	}
}

func newReit(ticker string, price float64) *resource.Security {
	return &resource.Security{
		Ticker:     ticker,
		Type:       resource.REIT_TYPE,
		Name:       "Name " + ticker,
		Admin:      "Admin " + ticker,
		Document:   "22.222.222/0001-22",
		Segment:    "Segment " + ticker,
		Currency:   brl,
		Price:      price,
		Origin:     "https://example.com/fundos-imobiliarios/" + ticker,
		CapturedAt: capturedAt,
	}
}

type fakeConfig struct{}

func (fakeConfig) GetVersion() string { return "development" }

func (fakeConfig) GetScrapingTimeout() time.Duration { return time.Second }

func (fakeConfig) GetMaxConcurrentRequests() int { return 4 }

type listFunc func(tickers []string) ([]*resource.Security, []*resource.TickerFailure)

// listing returns a listFunc that always answers with securities and no failures.
func listing(securities ...*resource.Security) listFunc {
	return func(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
		return securities, nil
	}
}

type fakeStockService struct {
	get  func(ticker string) *resource.Security
	list listFunc
}

func (f *fakeStockService) GetStockByTicker(ticker string) *resource.Security {
	return f.get(ticker)
}

func (f *fakeStockService) ListStocksByTickers(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
	return f.list(tickers)
}

type fakeReitService struct {
	get  func(ticker string) *resource.Security
	list listFunc
}

func (f *fakeReitService) GetReitByTicker(ticker string) *resource.Security {
	return f.get(ticker)
}

func (f *fakeReitService) ListReitsByTickers(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
	return f.list(tickers)
}

type purchaseCall struct {
	stocks []string
	reits  []string
	amount float64
}

type fakePurchaseBalanceService struct {
	result *resource.PurchaseBalance
	err    error
	calls  []purchaseCall
}

func (f *fakePurchaseBalanceService) PurchaseBalance(securities []*resource.Security, amountInvested float64) *resource.PurchaseBalance {
	return f.result
}

func (f *fakePurchaseBalanceService) PurchaseBalancesBySecurities(stockTickers []string, reitTickers []string, amountInvested float64) (*resource.PurchaseBalance, error) {
	f.calls = append(f.calls, purchaseCall{stocks: stockTickers, reits: reitTickers, amount: amountInvested})
	if f.err != nil {
		return nil, f.err
	}
	return f.result, nil
}

// fakes groups the fake services used to build a command tree.
type fakes struct {
	stock    *fakeStockService
	reit     *fakeReitService
	purchase *fakePurchaseBalanceService
}

func newFakes() *fakes {
	return &fakes{
		stock: &fakeStockService{
			get:  func(ticker string) *resource.Security { return nil },
			list: listing(),
		},
		reit: &fakeReitService{
			get:  func(ticker string) *resource.Security { return nil },
			list: listing(),
		},
		purchase: &fakePurchaseBalanceService{
			result: &resource.PurchaseBalance{},
		},
	}
}

// tree is a command tree built with fake services.
type tree struct {
	root   cmd.RootCommand
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func newTree(f *fakes) *tree {
	root := cmd.NewRootCommand(fakeConfig{})
	cmd.NewStockCommand(f.stock, f.purchase).InitApp(root)
	cmd.NewReitCommand(f.reit, f.purchase).InitApp(root)
	cmd.NewSecurityCommand(f.purchase).InitApp(root)
	return &tree{root: root, stdout: &bytes.Buffer{}, stderr: &bytes.Buffer{}}
}

// run executes the tree with args, collecting stdout and stderr in separate buffers.
// The buffers are reset on each call, so the same tree can be run more than once.
func (tr *tree) run(args ...string) error {
	tr.stdout.Reset()
	tr.stderr.Reset()
	c := tr.root.GetCobraCommand()
	c.SetOut(tr.stdout)
	c.SetErr(tr.stderr)
	c.SetArgs(args)
	return tr.root.Execute()
}
