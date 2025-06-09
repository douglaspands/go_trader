package service

import (
	"trader/internal/resource"
	"trader/internal/scraping"
)

func GetStock(ticker string) *resource.Security {
	stock, err := scraping.GetStockByTicker(ticker)
	if err != nil {
		return nil
	}
	return stock
}

func ListStocks(tickers []string) []*resource.Security {
	result := scraping.ListStocksByTickers(tickers)
	return result
}

func MakeStockPurchaseBalance(tickers []string, amountInvested float64) *resource.PurchaseBalance {
	stocks := ListStocks(tickers)
	result := MakePurchaseBalance(stocks, amountInvested)
	return result
}
