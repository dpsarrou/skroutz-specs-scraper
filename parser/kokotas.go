package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type KokotasParser struct{}

func (p *KokotasParser) ShopId() string {
	return "1283"
}

func (p *KokotasParser) ParseProductSpecs(html []byte) (map[string]string, error) {
	doc, err := parseHTML(html)
	if err != nil {
		return nil, err
	}
	data := map[string]string{}
	doc.Find("div.prod-main-info tr").Each(func(i int, s *goquery.Selection) {
		key := s.Find("th").Text()
		value := s.Find("td").Text()

		allowedValues := []string{
			"Καλύπτει χώρους (m²)",
			"Ονομαστική Κατανάλωση Ψύξη (kW)",
			"Ονομαστική Κατανάλωση Θέρμανση (kW)",
			"Φίλτρα",
			"Στάθμη Θορύβου Εσ. Μονάδας (dB(A))",
			"Στάθμη Θορύβου Εξ. Μονάδας (dB(A))",
			"Εγγύηση Προμηθευτή",
			"Χώρα Προέλευσης",
		}
		for _, v := range allowedValues {
			if strings.Contains(key, v) {
				data[key] = value
			}
		}
		if strings.Contains(key, "Πιστοποίηση Eurovent") {
			data[key] = "ΟΧΙ"
			if s.Find("td em").HasClass("yes") {
				data[key] = "ΝΑΙ"
			}
		}

		if strings.Contains(key, "Εύρος Απόδοσης - Ψύξη (Btu/h)") {
			value = strings.ReplaceAll(value, "~", "-")
			values := strings.Split(value, "-")
			if len(values) > 1 {
				data["Ελάχιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
				data["Μέγιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[1])
			} else {
				data["Ελάχιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
				data["Μέγιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
			}
		}

		if strings.Contains(key, "Εύρος Απόδοσης - Θέρμανση (Btu/h)") {
			value = strings.ReplaceAll(value, "~", "-")
			values := strings.Split(value, "-")
			if len(values) > 1 {
				data["Ελάχιστη Θέρμανση (Btu/h)"] = normalizeNumbers(values[0])
				data["Μέγιστη Θέρμανση (Btu/h)"] = values[1]
			} else {
				data["Ελάχιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
				data["Μέγιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
			}
		}

		if strings.Contains(key, "Ετήσια Κατανάλωση (kWh/a)") {
			data["Κατανάλωση Ψύξης"] = normalizeNumbers(value)
		}

		if strings.Contains(key, "Ετήσια Κατανάλωση - Μέση Ζώνη (kWh/a)") {
			data["Κατανάλωση Θέρμανσης"] = normalizeNumbers(value)
		}
	})
	return data, nil
}
