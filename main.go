package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"skroutz-specs-scraper/crawler"
	"skroutz-specs-scraper/exporter"
	"skroutz-specs-scraper/filter"
	"skroutz-specs-scraper/parser"

	"gopkg.in/yaml.v2"
)

func main() {

	var configFile = flag.String("config", "config.yml", "Filename of configuration file")
	var output = flag.String("output", "csv", "Output format")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Missing argument <url>")
	}
	url := flag.Arg(0)

	// Read configuration file
	config, err := loadConfigurationFile(*configFile)
	if err != nil {
		log.Fatalf("Could not load configuration file: %v", err)
	}

	// Initialize the crawler
	parsers := parser.LoadShopParsersByIds(config.ShopParsers)
	filters := filter.LoadFiltersById(config.Filters)
	crawler := crawler.NewCrawler(config.Crawler, parsers, filters)

	// Crawl products
	data, err := crawler.CrawlProductListPage(url)
	if err != nil {
		log.Fatal(err)
	}

	// Export to CSV
	if *output == "csv" {
		csv := exporter.CSV{Config: config.Export.CSV}
		output := csv.Export(data)
		fmt.Println(output)
	}
}

type Config struct {
	ShopParsers []string       `yaml:"shop_parsers"`
	Filters     []string       `yaml:"filters"`
	Crawler     crawler.Config `yaml:"crawler"`
	Export      struct {
		CSV exporter.CSVConfig `yaml:"csv"`
	} `yaml:"export"`
}

func loadConfigurationFile(filename string) (Config, error) {
	var config Config
	configData, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal([]byte(configData), &config)
	return config, err
}
