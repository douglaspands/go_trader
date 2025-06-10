package service

import (
	"math"
	"sort"
	"trader/internal/resource"
)

func MakePurchaseBalance(securities []*resource.Security, amountInvested float64) *resource.PurchaseBalance {
	securityCount := len(securities)
	securityValue := amountInvested / float64(securityCount)
	remainingBalance := amountInvested
	priceMin := math.MaxFloat64
	priceMax := float64(0)
	securitiesPurchase := make([]*resource.SecurityPurchase, 0)
	securitiesPurchaseSort := make([]*resource.SecurityPurchase, 0)
	securitiesExpensive := make([]*resource.Security, 0)
	countBalance := 0
	for _, security := range securities {
		if security.Price < priceMin {
			priceMin = security.Price
		}
		if security.Price > priceMax {
			priceMax = security.Price
		}
	}
	if amountInvested < priceMin {
		return &resource.PurchaseBalance{SecuritiesBalance: securitiesPurchase, AmountInvested: amountInvested}
	}
	if securityValue < priceMax {
		countBalance = -1
	}
	for _, security := range securities {
		count := int(securityValue / security.Price)
		if count < 1 {
			securitiesExpensive = append(securitiesExpensive, security)
		} else {
			if (count + countBalance) > 0 {
				count = count + countBalance
			}
			securityPurchase := &resource.SecurityPurchase{
				Security: security,
				Count:    count,
			}
			securitiesPurchase = append(securitiesPurchase, securityPurchase)
			securitiesPurchaseSort = append(securitiesPurchaseSort, securityPurchase)
			remainingBalance = remainingBalance - securityPurchase.TotalAmount()
		}
	}
	for _, security := range securitiesExpensive {
		if remainingBalance >= security.Price {
			count := 1
			securityPurchase := &resource.SecurityPurchase{
				Security: security,
				Count:    count,
			}
			securitiesPurchase = append(securitiesPurchase, securityPurchase)
			securitiesPurchaseSort = append(securitiesPurchaseSort, securityPurchase)
			remainingBalance = remainingBalance - securityPurchase.TotalAmount()
		}
	}
	sort.Slice(securitiesPurchaseSort, func(i, j int) bool {
		return securitiesPurchaseSort[i].Security.Price < securitiesPurchaseSort[j].Security.Price
	})
	for {
		if remainingBalance >= priceMin {
			for i := range securityCount {
				if remainingBalance >= securitiesPurchaseSort[i].Security.Price {
					securitiesPurchaseSort[i].Count += 1
					remainingBalance = remainingBalance - securitiesPurchaseSort[i].Security.Price
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
