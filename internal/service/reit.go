package service

import (
	"trader/internal/resource"
	"trader/internal/scraping"
)

func GetReit(ticker string) *resource.Security {
	stock, err := scraping.GetReitByTicker(ticker)
	if err != nil {
		return nil
	}
	return stock
}

func ListReits(tickers []string) []*resource.Security {
	result := scraping.ListReitsByTickers(tickers)
	return result
}

func MakeReitPurchaseBalance(tickers []string, amountInvested float64) *resource.PurchaseBalance {
	reits := ListReits(tickers)
	result := MakePurchaseBalance(reits, amountInvested)
	return result
}
