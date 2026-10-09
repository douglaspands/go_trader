package service

import (
	"trader/internal/resource"
	"trader/internal/scraping"
)

type ReitService interface {
	GetReitByTicker(ticker string) *resource.Security
	ListReitsByTickers(tickers []string) []*resource.Security
}

type reitService struct {
	reitScraping scraping.ReitScraping
}

func (rs *reitService) GetReitByTicker(ticker string) *resource.Security {
	reit, err := rs.reitScraping.GetReitByTicker(ticker)
	if err != nil {
		return nil
	}
	return reit
}

func (rs *reitService) ListReitsByTickers(tickers []string) []*resource.Security {
	result := rs.reitScraping.ListReitsByTickers(tickers)
	return result
}

func NewReitService(reitScraping scraping.ReitScraping) ReitService {
	return &reitService{
		reitScraping: reitScraping,
	}
}
