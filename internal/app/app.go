package app

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
	cmdRoot cmd.RootCommand
}

func (a *app) setup() {
	config := config.NewConfig()
	stockScraping := scraping.NewStockScraping(config)
	stockService := service.NewStockService(stockScraping)
	reitScraping := scraping.NewReitScraping(config)
	reitService := service.NewReitService(reitScraping)
	purchaseBalanceService := service.NewPurchaseBalanceService(stockService, reitService)

	a.cmdRoot = cmd.NewRootCommand(config)
	stockCommand := cmd.NewStockCommand(stockService)
	stockCommand.InitApp(a.cmdRoot)

	reitCommand := cmd.NewReitCommand(reitService)
	reitCommand.InitApp(a.cmdRoot)

	securityCommand := cmd.NewSecurityCommand(purchaseBalanceService)
	securityCommand.InitApp(a.cmdRoot)

}

func (a *app) Run() {
	a.setup()
	err := a.cmdRoot.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func NewApp() App {
	return &app{}
}
