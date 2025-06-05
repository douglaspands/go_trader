package scraping

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
	"trader/internal/resource"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func GetStockByTicker(ticker string) (*resource.Stock, error) {

	url := fmt.Sprintf("%s/acoes/%s", STOCK_HOST_URL, strings.ToLower(ticker))
	htmlDoc, err := getHtml(url)
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
		name = strings.TrimSpace(strings.Split(htmlquery.SelectAttr(n, "title"), "-")[1])
	}

	n = htmlquery.FindOne(doc, `//div[@title="Valor atual do ativo"]/strong/text()`)
	var price float64
	if n != nil {
		price, err = strconv.ParseFloat(strings.ReplaceAll(n.Data, ",", "."), 64)
		if err != nil {
			return nil, err
		}
	}

	n = htmlquery.FindOne(doc, `//*[@id='company-section']/div[1]/div/div[1]/div[2]/h4/small/text()`)
	var document string
	if n != nil {
		document = strings.TrimSpace(n.Data)
	}

	var description []string
	for _, n = range htmlquery.Find(doc, `//div/p[not(@*)]/text()`) {
		description = append(description, strings.TrimSpace(n.Data))
	}

	return &resource.Stock{
		Ticker:      strings.ToUpper(ticker),
		Name:        name,
		Description: strings.TrimSpace(strings.Join(description, " ")),
		Price:       price,
		Document:    document,
		Origin:      url,
		CapturedAt:  time.Now(),
	}, nil
}

func ListStocksByTickers(tickers []string) []*resource.Stock {
	var stocks []*resource.Stock
	for _, ticker := range tickers {
		stock, err := GetStockByTicker(ticker)
		if err == nil {
			stocks = append(stocks, stock)
		}
	}
	return stocks
}
