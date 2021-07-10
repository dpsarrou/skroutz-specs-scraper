package exporter

import (
	"strings"
)

type CSV struct {
	Config CSVConfig
}

type CSVConfig struct {
	AllHeaders bool     `yaml:"include_all_headers"`
	Headers    []string `yaml:"headers"`
}

func NewCsvExporter(config CSVConfig) *CSV {
	return &CSV{config}
}

func (c *CSV) Export(data []map[string]string) string {
	rows := []string{}
	if len(data) > 0 {

		headers := extractHeaders(data, c.Config.Headers, c.Config.AllHeaders)
		rows = append(rows, strings.Join(headers, "\t"))

		for _, product := range data {
			row := []string{}
			for _, header := range headers {
				value, exists := product[header]
				if exists != true {
					value = ""
				}
				row = append(row, value)
			}
			rows = append(rows, strings.Join(row, "\t"))
		}
	}
	return strings.Join(rows, "\n")
}

func extractHeaders(data []map[string]string, headers []string, allHeaders bool) []string {
	if allHeaders {
		headersList := map[string]int{}
		for _, product := range data {
			for key := range product {
				headersList[key] = 1
			}
		}

		headers = sortHeaders(headers, headersList)
	}
	return headers
}

func sortHeaders(headers []string, headerList map[string]int) []string {
	existingHeadersMap := map[string]int{}
	for _, header := range headers {
		existingHeadersMap[header] = 1
	}

	for header := range headerList {
		_, exists := existingHeadersMap[header]
		if exists {
			continue
		}
		headers = append(headers, header)
	}

	return headers
}
