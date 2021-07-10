package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SkroutzParser struct{}

type Product struct {
	Specs     map[string]string
	StoreUrls map[string]string
}

func (p *SkroutzParser) ParseProductSinglePage(html []byte) (*Product, error) {
	doc, err := parseHTML(html)
	if err != nil {
		return nil, err
	}
	specs := p.parseProductSpecs(doc)
	storeUrls := p.parseShopProductUrls(doc)

	product := &Product{
		Specs:     specs,
		StoreUrls: storeUrls,
	}
	return product, nil
}

func (p *SkroutzParser) parseProductSpecs(doc *goquery.Document) map[string]string {
	data := map[string]string{}
	// Url
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("rel")
		if exists && value == "canonical" {
			data["Url"], _ = s.Attr("href")
		}
	})

	// Brand
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("itemprop")
		if exists && value == "brand" {
			data["Brand"], _ = s.Attr("content")
		}
	})

	// Product
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		product := ""
		value, exists := s.Attr("property")
		if exists && value == "og:title" {
			product, _ = s.Attr("content")
			brand, exists := data["Brand"]
			if exists {
				product = strings.ReplaceAll(product, brand+" ", "")
			}
			data["Product"] = product
		}
	})

	// Reviews
	doc.Find("section.sku-details div.rating-wrapper div.actual-rating").Each(func(i int, s *goquery.Selection) {
		value := s.Text()
		if value == "" {
			value = "0"
		}
		data["Αξιολογήσεις"] = value
	})

	// Rating
	doc.Find("section.sku-details div.rating-wrapper span").Each(func(i int, s *goquery.Selection) {
		value := s.Text()
		if value == "0.0" {
			value = ""
		}
		data["Βαθμολογία"] = value
	})

	// Price
	doc.Find("li.card.card-variant.selected span.price-details em").Each(func(i int, s *goquery.Selection) {
		price := strings.ReplaceAll(s.Text(), " €", "")
		price = normalizeNumbers(price)
		data["Τιμή"] = price
	})

	// Specs
	doc.Find("div.specs-container > div.spec-groups div.spec-details dl").Each(func(i int, dl *goquery.Selection) {
		key := dl.Find("dt").Text()
		value := dl.Find("dd").Text()
		data[key] = value
	})

	return data
}

func (p *SkroutzParser) parseShopProductUrls(doc *goquery.Document) map[string]string {
	data := map[string]string{}
	doc.Find("div.prices ol.sku-list li.card").Each(func(i int, s *goquery.Selection) {
		store, _ := s.Attr("id")
		store = strings.ReplaceAll(store, "shop-", "")
		storeUrl := ""
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			storeUrl, _ = s.Attr("href")
		})
		data[store] = "https://www.skroutz.gr" + storeUrl
	})
	return data
}

func (p *SkroutzParser) ParseShopRedirectUrl(html []byte) (string, error) {
	url := ""
	doc, err := parseHTML(html)
	if err != nil {
		return url, err
	}
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("rel")
		if exists && value == "prefetch" {
			url, _ = s.Attr("href")
		}
	})
	return url, nil
}

func (p *SkroutzParser) ParseProductsListPage(html []byte) ([]string, string, error) {
	urls := []string{}
	hasNext := ""
	doc, err := parseHTML(html)
	if err != nil {
		return urls, hasNext, err
	}
	doc.Find("ol.list li.card h2 a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		urls = append(urls, "https://www.skroutz.gr"+href)
	})

	doc.Find("ol.paginator li").Last().Find("a").Each(func(i int, s *goquery.Selection) {
		hasNext, _ = s.Attr("href")
		if hasNext != "" {
			hasNext = "https://www.skroutz.gr" + hasNext
		}
	})
	return urls, hasNext, nil
}
