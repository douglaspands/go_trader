package core

import (
	"os"
	"trader/cmd"
	"trader/internal/config"
	"trader/internal/scraping"
	"trader/internal/service"
)

type App interface {
	Run()
}

type app struct {
	config                 config.Config
	stockScraping          scraping.StockScraping
	reitScraping           scraping.ReitScraping
	stockService           service.StockService
	reitService            service.ReitService
	purchaseBalanceService service.PurchaseBalanceService
	// Commands
	rootCommand     cmd.RootCommand
	stockCommand    cmd.StockCommand
	reitCommand     cmd.ReitCommand
	securityCommand cmd.SecurityCommand
}

func (a *app) setup() {
	a.config = config.NewConfig()
	a.stockScraping = scraping.NewStockScraping(a.config)
	a.reitScraping = scraping.NewReitScraping(a.config)
	a.stockService = service.NewStockService(a.stockScraping)
	a.reitService = service.NewReitService(a.reitScraping)
	a.purchaseBalanceService = service.NewPurchaseBalanceService(a.stockService, a.reitService)

	a.rootCommand = cmd.NewRootCommand(a.config)
	a.stockCommand = cmd.NewStockCommand(a.stockService, a.purchaseBalanceService)
	a.stockCommand.InitApp(a.rootCommand)

	a.reitCommand = cmd.NewReitCommand(a.reitService, a.purchaseBalanceService)
	a.reitCommand.InitApp(a.rootCommand)

	a.securityCommand = cmd.NewSecurityCommand(a.purchaseBalanceService)
	a.securityCommand.InitApp(a.rootCommand)

}

func (a *app) Run() {
	a.setup()
	err := a.rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func NewApp() App {
	return &app{}
}
