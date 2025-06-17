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
	rootCommand cmd.RootCommand
}

func (a *app) setup() {
	config := config.NewConfig()

	stockScraping := scraping.NewStockScraping(config)
	reitScraping := scraping.NewReitScraping(config)

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
