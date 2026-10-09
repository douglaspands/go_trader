package service

import (
	"math"
	"sort"
	"trader/internal/resource"
)

type PurchaseBalanceService interface {
	PurchaseBalance(securities []*resource.Security, amountInvested float64) *resource.PurchaseBalance
	PurchaseBalancesBySecurities(stockTickers []string, reitTickers []string, amountInvested float64) *resource.PurchaseBalance
}

type purchaseBalanceService struct {
	stockService StockService
	reitService  ReitService
}

func (pb *purchaseBalanceService) PurchaseBalance(allSecurities []*resource.Security, amountInvested float64) *resource.PurchaseBalance {
	securities := make([]*resource.Security, 0, len(allSecurities))
	for _, security := range allSecurities {
		if security != nil && security.Price > 0 {
			securities = append(securities, security)
		}
	}
	securityCount := len(securities)
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
	if securityCount == 0 || amountInvested < priceMin {
		return &resource.PurchaseBalance{SecuritiesBalance: securitiesPurchase, AmountInvested: amountInvested}
	}
	securityValue := amountInvested / float64(securityCount)
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
	for len(securitiesPurchaseSort) > 0 && remainingBalance >= securitiesPurchaseSort[0].Security.Price {
		for _, securityPurchase := range securitiesPurchaseSort {
			if remainingBalance >= securityPurchase.Security.Price {
				securityPurchase.Count += 1
				remainingBalance = remainingBalance - securityPurchase.Security.Price
			}
		}
	}
	return &resource.PurchaseBalance{SecuritiesBalance: securitiesPurchase, AmountInvested: amountInvested}
}

func (pb *purchaseBalanceService) PurchaseBalancesBySecurities(stockTickers []string, reitTickers []string, amountInvested float64) *resource.PurchaseBalance {
	stocks := pb.stockService.ListStocksByTickers(stockTickers)
	reits := pb.reitService.ListReitsByTickers(reitTickers)
	securities := append(stocks, reits...)
	result := pb.PurchaseBalance(securities, amountInvested)
	return result
}

func NewPurchaseBalanceService(stockService StockService, reitService ReitService) PurchaseBalanceService {
	return &purchaseBalanceService{
		stockService: stockService,
		reitService:  reitService,
	}
}
