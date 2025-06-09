package service

import (
	"math"
	"trader/internal/resource"
)

func MakePurchaseBalance(securities []*resource.Security, amountInvested float64) *resource.PurchaseBalance {
	stockCount := len(securities)
	stockValue := amountInvested / float64(stockCount)
	remainingBalance := amountInvested
	priceMin := math.MaxFloat64
	securitiesPurchase := make([]*resource.SecurityPurchase, 0)
	for _, stock := range securities {
		if stock.Price < priceMin {
			priceMin = stock.Price
		}
		securityPurchase := &resource.SecurityPurchase{
			Security: stock,
			Count:    int(stockValue / stock.Price),
		}
		securitiesPurchase = append(securitiesPurchase, securityPurchase)
		remainingBalance = remainingBalance - securityPurchase.TotalAmount()
	}
	for {
		if remainingBalance >= priceMin {
			for i := range stockCount {
				if remainingBalance >= securitiesPurchase[i].Security.Price {
					securitiesPurchase[i].Count += 1
					remainingBalance = remainingBalance - securitiesPurchase[i].Security.Price
				}
			}
		} else {
			break
		}
	}
	return &resource.PurchaseBalance{SecuritiesBalance: securitiesPurchase, AmountInvested: amountInvested}
}

func MakeSecuritiesPurchaseBalance(stockTickers []string, reitTickers []string, amountInvested float64) *resource.PurchaseBalance {
	stocks := ListStocks(stockTickers)
	reits := ListReits(reitTickers)
	securities := append(stocks, reits...)
	result := MakePurchaseBalance(securities, amountInvested)
	return result
}
