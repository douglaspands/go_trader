package scraping

import (
	"bytes"
	"fmt"
	"strings"
	"time"
	"trader/internal/config"
	"trader/internal/resource"
	"trader/internal/tools"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type StockScraping interface {
	GetStockByTicker(ticker string) (*resource.Security, error)
	ListStocksByTickers(tickers []string) ([]*resource.Security, []*resource.TickerFailure)
}

type stockScraping struct {
	url     string
	config  config.Config
	fetcher *Fetcher
}

func (ss *stockScraping) GetStockByTicker(ticker string) (*resource.Security, error) {
	if err := validateTicker(ticker); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/acoes/%s", ss.url, strings.ToLower(ticker))
	htmlDoc, err := ss.fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(bytes.NewReader(htmlDoc))
	if err != nil {
		return nil, err
	}

	var n *html.Node
	n = htmlquery.FindOne(doc, "//h1[@title]")
	var name string
	if n != nil {
		if parts := strings.SplitN(htmlquery.SelectAttr(n, "title"), "-", 2); len(parts) == 2 {
			name = strings.TrimSpace(parts[1])
		}
	}

	n = htmlquery.FindOne(doc, `//div[@title="Valor atual do ativo"]/strong/text()`)
	var price float64 = 0.0
	if n != nil {
		price = tools.ToFloat(strings.TrimSpace(n.Data), ",")
	}

	n = htmlquery.FindOne(doc, `//*[@id='company-section']/div[1]/div/div[1]/div[2]/h4/small/text()`)
	var document string
	if n != nil {
		document = strings.TrimSpace(n.Data)
	}

	return &resource.Security{
		Ticker: strings.ToUpper(ticker),
		Name:   name,
		Type:   resource.STOCK_TYPE,
		Currency: &resource.Currency{
			Code:        "BRL",
			Description: "Brazilian Real",
			Sign:        "R$",
		},
		Price:      price,
		Document:   document,
		Origin:     url,
		CapturedAt: time.Now(),
	}, nil
}

func (ss *stockScraping) ListStocksByTickers(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
	return listByTickers(tickers, resource.STOCK_TYPE, ss.GetStockByTicker)
}

func NewStockScraping(config config.Config, fetcher *Fetcher) StockScraping {
	return &stockScraping{
		url:     STATUS_INVEST_URL,
		config:  config,
		fetcher: fetcher,
	}
}
