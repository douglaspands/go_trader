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

func GetReitByTicker(ticker string) (*resource.Security, error) {

	url := fmt.Sprintf("%s/fundos-imobiliarios/%s", STATUS_INVEST_URL, strings.ToLower(ticker))
	htmlDoc, err := getHtml(url)
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(bytes.NewReader(htmlDoc))
	if err != nil {
		return nil, err
	}

	var n *html.Node
	n = htmlquery.FindOne(doc, "//h1[@class='lh-4']/small/text()")
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
	var price float64
	if n != nil {
		price, err = strconv.ParseFloat(strings.ReplaceAll(n.Data, ",", "."), 64)
		if err != nil {
			return nil, err
		}
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

func ListReitsByTickers(tickers []string) []*resource.Security {
	var reits []*resource.Security
	for _, ticker := range tickers {
		reit, err := GetReitByTicker(ticker)
		if err == nil {
			reits = append(reits, reit)
		}
	}
	return reits
}
