package parser

import (
	"bytes"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ShopParser interface {
	ShopId() string
	ParseProductSpecs(html []byte) (map[string]string, error)
}

var allParsers []ShopParser = []ShopParser{
	&KokotasParser{},
	&AbClimaParser{},
	// Add new parser below
}

func LoadShopParsersByIds(ids []string) []ShopParser {
	var parsers []ShopParser
	for _, id := range ids {
		for _, parser := range allParsers {
			if id == parser.ShopId() {
				parsers = append(parsers, parser)
			}
		}
	}
	return parsers
}

func normalizeNumbers(number string) string {
	n := strings.ReplaceAll(number, ".", "")
	n = strings.ReplaceAll(n, ",", ".")
	return n
}

func parseHTML(html []byte) (*goquery.Document, error) {
	reader := io.NopCloser(bytes.NewBuffer(html))
	doc, err := goquery.NewDocumentFromReader(reader)
	return doc, err
}
