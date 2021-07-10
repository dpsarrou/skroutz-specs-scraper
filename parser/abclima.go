package parser

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AbClimaParser struct{}

func (p *AbClimaParser) ShopId() string {
	return "3292"
}

func (p *AbClimaParser) ParseProductSpecs(html []byte) (map[string]string, error) {
	doc, err := parseHTML(html)
	if err != nil {
		return nil, err
	}
	data := map[string]string{}
	doc.Find("div.product_extra li").Each(func(i int, s *goquery.Selection) {

		if s.Has("table").Length() > 0 {

			s.Find("table tr").Each(func(i int, s *goquery.Selection) {
				data = p.parseSpecs(data, s.Find("td span"))
			})

		} else if s.Has("li").Length() > 0 {

			log.Println("Skipping product")
		} else {
			data = p.parseSpecs(data, s)
		}

	})
	return data, nil
}

func (p *AbClimaParser) parseSpecs(data map[string]string, s *goquery.Selection) map[string]string {

	property := s.Text()

	if strings.Contains(property, "Ελάχιστη απόδοση:") {
		values := strings.Split(strings.ReplaceAll(property, "Ελάχιστη απόδοση: ", ""), ",")
		data["Ελάχιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
		data["Ελάχιστη Θέρμανση (Btu/h)"] = normalizeNumbers(values[1])
	}

	if strings.Contains(property, "Μέγιστη απόδοση:") {
		values := strings.Split(strings.ReplaceAll(property, "Μέγιστη απόδοση: ", ""), ",")
		data["Μέγιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
		data["Μέγιστη Θέρμανση (Btu/h)"] = normalizeNumbers(values[1])
	}

	// Noise
	if strings.Contains(property, "Επίπεδο Θορύβου") {
		if strings.Contains(property, "Επίπεδο Θορύβου dB(A)Lp") {
			if strings.Contains(property, "Εσωτερικής Μονάδας") {
				values := strings.Split(property, ":")
				data["Στάθμη Θορύβου Εσ. Μονάδας (dB(A))"] = strings.ReplaceAll(values[1], " ", "")
			}
			if strings.Contains(property, "Εξωτερικής Μονάδας") {
				values := strings.Split(property, ":")
				data["Στάθμη Θορύβου Εξ. Μονάδας (dB(A))"] = strings.ReplaceAll(values[1], " ", "")
			}
		} else if strings.Contains(property, "Επίπεδο Θορύβου (dB(A)±3)") {
			values := strings.Split(property, ":")
			if strings.Contains(values[1], ",") {
				noise := strings.Split(values[1], ",")
				data["Στάθμη Θορύβου Εσ. Μονάδας (dB(A))"] = noise[0]
				data["Στάθμη Θορύβου Εξ. Μονάδας (dB(A))"] = noise[1]
			} else {
				data["Στάθμη Θορύβου Εσ. Μονάδας (dB(A))"] = values[1]
			}

		} else if strings.Contains(property, "Επίπεδο Θορύβου σε db") {
			if strings.Contains(property, "Εσωτερικής Μονάδας") {
				values := strings.Split(property, "Εσωτερικής Μονάδας")
				data["Στάθμη Θορύβου Εσ. Μονάδας (dB(A))"] = strings.ReplaceAll(values[1], " ", "")
			}
			if strings.Contains(property, "Εξωτερικής Μονάδας") {
				values := strings.Split(property, "Εξωτερικής Μονάδας")
				data["Στάθμη Θορύβου Εξ. Μονάδας (dB(A))"] = strings.ReplaceAll(values[1], " ", "")
			}
		} else {
			log.Fatal(property)
		}
	} else if strings.Contains(property, "Ηχητική Πίεση Εσωτερικής Μονάδας (dB(Α))") {
		noise := strings.Split(property, "Ηχητική Πίεση Εσωτερικής Μονάδας (dB(Α))")
		data["Στάθμη Θορύβου Εσ. Μονάδας (dB(A))"] = noise[1]
	} else if strings.Contains(property, "Ηχητική Πίεση Εξωτερικής Μονάδας (dB(Α))") {
		noise := strings.Split(property, "Ηχητική Πίεση Εξωτερικής Μονάδας (dB(Α))")
		data["Στάθμη Θορύβου Εξ. Μονάδας (dB(A))"] = noise[1]
	}

	space := "συνιστάται για χώρους εως m2"
	key := strings.ToLower(property)
	if strings.Contains(key, space) {
		data["Καλύπτει χώρους (m²)"] = strings.ReplaceAll(key, space, "")
	} else if strings.Contains(key, "καλύπτει χώρους (m²)") {
		values := strings.Split(property, "Καλύπτει χώρους (m²)")
		data["Καλύπτει χώρους (m²)"] = values[1]
	}

	if strings.Contains(s.Text(), "Εγγύηση") {
		data["Εγγύηση Προμηθευτή"] = strings.ReplaceAll(s.Text(), "Εγγύηση", "")
	}

	if strings.Contains(property, "Εύρος Απόδοσης - Ψύξη (Btu/h)") {
		split := strings.Split(property, "Εύρος Απόδοσης - Ψύξη (Btu/h)")
		value := strings.ReplaceAll(split[1], "~", "-")
		value = strings.ReplaceAll(value, " ", "")
		values := strings.Split(value, "-")
		if len(values) > 1 {
			data["Ελάχιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
			data["Μέγιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[1])
		} else {
			data["Ελάχιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
			data["Μέγιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
		}
	}

	if strings.Contains(property, "Εύρος Απόδοσης - Θέρμανση (Btu/h)") {
		split := strings.Split(property, "Εύρος Απόδοσης - Θέρμανση (Btu/h)")
		value := strings.ReplaceAll(split[1], "~", "-")
		value = strings.ReplaceAll(value, " ", "")
		values := strings.Split(value, "-")
		if len(values) > 1 {
			data["Ελάχιστη Θέρμανση (Btu/h)"] = normalizeNumbers(values[0])
			data["Μέγιστη Θέρμανση (Btu/h)"] = normalizeNumbers(values[1])
		} else {
			data["Ελάχιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
			data["Μέγιστη Ψύξη (Btu/h)"] = normalizeNumbers(values[0])
		}
	}

	if strings.Contains(property, "Πιστοποίηση Eurovent") {
		values := strings.Split(property, "Πιστοποίηση Eurovent")
		data["Πιστοποίηση Eurovent"] = values[1]
	}

	return data
}
