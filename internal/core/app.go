package core

import (
	"trader/cmd"
	"trader/internal/config"
	"trader/internal/scraping"
	"trader/internal/service"
)

type App interface {
	Run(args []string) int
}

type app struct {
	rootCommand cmd.RootCommand
}

func (a *app) setup() {
	config := config.NewConfig()

	fetcher := scraping.NewFetcher(config)
	stockScraping := scraping.NewStockScraping(config, fetcher)
	reitScraping := scraping.NewReitScraping(config, fetcher)

	stockService := service.NewStockService(stockScraping)
	reitService := service.NewReitService(reitScraping)
	purchaseBalanceService := service.NewPurchaseBalanceService(stockService, reitService)

	rootCommand := cmd.NewRootCommand(config)

	stockCommand := cmd.NewStockCommand(stockService, purchaseBalanceService)
	stockCommand.InitApp(rootCommand)

	reitCommand := cmd.NewReitCommand(reitService, purchaseBalanceService)
	reitCommand.InitApp(rootCommand)

	securityCommand := cmd.NewSecurityCommand(purchaseBalanceService)
	securityCommand.InitApp(rootCommand)

	a.rootCommand = rootCommand
}

func (a *app) Run(args []string) int {
	a.setup()
	a.rootCommand.GetCobraCommand().SetArgs(args)
	if err := a.rootCommand.Execute(); err != nil {
		return 1
	}
	return 0
}

func NewApp() App {
	return &app{}
}
