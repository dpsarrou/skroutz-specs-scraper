package crawler

import (
	"log"
	"skroutz-specs-scraper/filter"
	"skroutz-specs-scraper/parser"
)

type SkroutzCrawler struct {
	SkroutzClient HttpClient
	ShopClient    HttpClient
	Parser        *parser.SkroutzParser
	ShopParsers   []parser.ShopParser
	Filters       []filter.Filter
}

type Config struct {
	UserAgent     string `yaml:"user_agent"`
	SkroutzCookie string `yaml:"skroutz_cookie"`
}

func NewCrawler(config Config, parsers []parser.ShopParser, filters []filter.Filter) *SkroutzCrawler {
	skroutzClient := NewClient(config.UserAgent, config.SkroutzCookie)
	shopClient := NewClient(config.UserAgent, "")
	p := &parser.SkroutzParser{}
	return &SkroutzCrawler{
		SkroutzClient: skroutzClient, ShopClient: shopClient, Parser: p, ShopParsers: parsers, Filters: filters,
	}
}

func (c *SkroutzCrawler) CrawlProductListPage(url string) ([]map[string]string, error) {
	urls, err := c.getAllProductUrls(url)
	if err != nil {
		return nil, err
	}
	var productData []map[string]string
	total := len(urls)
	for index, url := range urls {

		log.Printf("[Crawler]\tCrawling product %d of %d..", index+1, total)
		productPage, err := c.SkroutzClient.Get(url)
		if err != nil {
			return nil, err
		}

		product, err := c.Parser.ParseProductSinglePage(productPage)
		if err != nil {
			return nil, err
		}

		var storeProductData map[string]string
		for _, parser := range c.ShopParsers {
			id := parser.ShopId()
			if storeRedirectUrl, ok := product.StoreUrls[id]; ok {
				storeProductPage, err := c.followStoreProductPage(storeRedirectUrl)
				if err != nil {
					return nil, err
				}
				storeProductData, err = parser.ParseProductSpecs(storeProductPage)
				if err != nil {
					return nil, err
				}
				// Use only one shop to parse data
				break
			}
		}

		for key, value := range storeProductData {
			product.Specs[key] = value
		}

		// Filter data
		for _, filter := range c.Filters {
			product.Specs = filter.Filter(product.Specs)
		}

		productData = append(productData, product.Specs)

	}
	return productData, nil
}

func (c *SkroutzCrawler) followStoreProductPage(url string) ([]byte, error) {
	storeRedirectPage, err := c.SkroutzClient.Get(url)
	if err != nil {
		return nil, err
	}
	storeProductUrl, err := c.Parser.ParseShopRedirectUrl(storeRedirectPage)
	if err != nil {
		return nil, err
	}
	storeProductPage, err := c.ShopClient.Get(storeProductUrl)
	return storeProductPage, err
}

func (c *SkroutzCrawler) getAllProductUrls(url string) ([]string, error) {
	var urls []string
	for {
		html, err := c.SkroutzClient.Get(url)
		if err != nil {
			return urls, err
		}
		productUrls, next, err := c.Parser.ParseProductsListPage(html)
		if err != nil {
			return urls, err
		}
		urls = append(urls, productUrls...)
		if next == "" {
			break
		}
		url = next
	}
	return urls, nil
}
