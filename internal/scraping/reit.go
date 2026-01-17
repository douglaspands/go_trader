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

const (
	REIT_TICKER_NOT_FOUND_MESSAGE = "Não encontramos o que você está procurando"
)

type ReitScraping interface {
	GetReitByTicker(ticker string) (*resource.Security, error)
	ListReitsByTickers(tickers []string) []*resource.Security
}
type reitScraping struct {
	url    string
	config config.Config
}

func (rs *reitScraping) GetReitByTicker(ticker string) (*resource.Security, error) {

	url := fmt.Sprintf("%s/fundos-imobiliarios/%s", STATUS_INVEST_URL, strings.ToLower(ticker))
	timeout := rs.config.GetScrapingTimeout()
	htmlDoc, _ := getHtml(url, timeout)

	doc, _ := htmlquery.Parse(bytes.NewReader(htmlDoc))

	ns := htmlquery.Find(doc, `//*[@id="main-2"]/section/div/h1/text()`)
	for _, n := range ns {
		if strings.TrimSpace(n.Data) == REIT_TICKER_NOT_FOUND_MESSAGE {
			return nil, fmt.Errorf("ticker not found for url=\"%s\"", url)
		}
	}

	var n *html.Node
	n = htmlquery.FindOne(doc, `//h1[@class='lh-4']/small/text()`)
	var name string
	if n != nil {
		name = strings.TrimSpace(n.Data)
	}

	n = htmlquery.FindOne(doc, `//*[@id='fund-section']/div/div/div[2]/div/div[1]/div/div/strong/text()`)
	var document string
	if n != nil {
		document = strings.TrimSpace(n.Data)
	}

	n = htmlquery.FindOne(doc, `//*[@id='fund-section']/div/div/div[2]/div/div[6]/div/div/strong/text()`)
	var segment string
	if n != nil {
		segment = strings.TrimSpace(n.Data)
	}

	n = htmlquery.FindOne(doc, `//div[@title="Valor atual do ativo"]/strong/text()`)
	var price float64 = 0.0
	if n != nil {
		price = tools.ToFloat(n.Data, ",")
	}

	n = htmlquery.FindOne(doc, `//*[@id='fund-section']/div/div/div[3]/div/div[2]/div[1]/div/strong/text()`)
	var admin string
	if n != nil {
		admin = strings.TrimSpace(n.Data)
	}

	return &resource.Security{
		Ticker:  strings.ToUpper(ticker),
		Name:    name,
		Admin:   admin,
		Segment: segment,
		Type:    resource.REIT_TYPE,
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

func (rs *reitScraping) ListReitsByTickers(tickers []string) []*resource.Security {
	var reits []*resource.Security
	for _, ticker := range tickers {
		reit, err := rs.GetReitByTicker(ticker)
		if err == nil {
			reits = append(reits, reit)
		}
	}
	return reits
}

func NewReitScraping(config config.Config) ReitScraping {
	return &reitScraping{
		url:    STATUS_INVEST_URL,
		config: config,
	}
}
