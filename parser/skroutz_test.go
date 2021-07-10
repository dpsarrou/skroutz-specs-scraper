package parser

import (
	"log"
	"os"
	"testing"
)

func readFile(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func TestItParsesSampleProductFile(t *testing.T) {

	d := readFile("../testdata/skroutz_single_product_page.html")
	p := &SkroutzParser{}
	product, err := p.ParseProductSinglePage(d)
	if err != nil {
		t.Error(err)
	}

	if product.Specs["Brand"] != "Tesla" {
		t.Errorf("Could not parse Brand")
	}

	if product.Specs["Τιμή"] != "638.49" {
		t.Errorf("Could not parse Price")
	}

	if product.Specs["Αξιολογήσεις"] != "25" {
		t.Errorf("Could not parse Αξιολογήσεις")
	}

	if product.Specs["Βαθμολογία"] != "4.9" {
		t.Errorf("Could not parse Βαθμολογία")
	}
}

func TestItParsesSampleFileWithoutRating(t *testing.T) {
	d := readFile("../testdata/skroutz_single_product_page_without_rating.html")
	p := &SkroutzParser{}
	product, err := p.ParseProductSinglePage(d)
	if err != nil {
		t.Error(err)
	}

	if product.Specs["Αξιολογήσεις"] != "0" {
		t.Errorf("Could not parse Αξιολογήσεις")
	}

	if product.Specs["Βαθμολογία"] != "" {
		t.Errorf("Could not parse Βαθμολογία")
	}
}

func TestItParsesShopProductUrls(t *testing.T) {
	d := readFile("../testdata/skroutz_single_product_page.html")
	p := &SkroutzParser{}
	product, err := p.ParseProductSinglePage(d)
	if err != nil {
		t.Error(err)
	}
	if product.StoreUrls["2843"] != "https://www.skroutz.gr/products/show/59670135" {
		t.Errorf("Could not parse shop product urls")
	}
}

func TestItParsesShopRedirectUrl(t *testing.T) {
	d := readFile("../testdata/shop_redirect.html")
	p := &SkroutzParser{}
	url, err := p.ParseShopRedirectUrl(d)
	if err != nil {
		t.Error(err)
	}
	if url != "https://abclima.gr/index.php?route=product/product&path=1_15_21_22&product_id=12995" {
		t.Errorf("Could not parse shop redirect url")
	}
}

func TestItCrawlsSkroutzProductUrls(t *testing.T) {
	d := readFile("../testdata/skroutz_products_1.html")
	p := &SkroutzParser{}
	data, _, err := p.ParseProductsListPage(d)
	if err != nil {
		t.Error(err)
	}
	if len(data) == 0 {
		t.Errorf("Could not parse Skroutz product urls")
	}
}

func TestItCrawlsSkroutzProductUrlsAndFindsSingleUrl(t *testing.T) {
	d := readFile("../testdata/skroutz_products_list_just_1_item.html")
	p := &SkroutzParser{}
	data, _, err := p.ParseProductsListPage(d)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		log.Println(data)
		t.Errorf("Expecting 1 url, found %d", len(data))
	}
}

func TestItParsesSkroutzPaginationProperly(t *testing.T) {
	d := readFile("../testdata/skroutz_products_1.html")
	p := &SkroutzParser{}
	_, hasNext, err := p.ParseProductsListPage(d)
	if err != nil {
		t.Error(err)
	}

	if hasNext == "" {
		t.Errorf("Expecting next page url, found none")
	}

	d = readFile("../testdata/skroutz_products_2.html")
	_, hasNext, err = p.ParseProductsListPage(d)
	if err != nil {
		t.Error(err)
	}

	if hasNext == "" {
		t.Errorf("Expecting next page url, found none")
	}

	d = readFile("../testdata/skroutz_products_3.html")
	_, hasNext, err = p.ParseProductsListPage(d)
	if err != nil {
		t.Error(err)
	}

	if hasNext != "" {
		t.Errorf("Found pagination url while expecting none")
	}
}
