package filter

import "strings"

type AirconditionSpecsNormalizer struct{}

func (n *AirconditionSpecsNormalizer) Id() string {
	return "AirconditionSpecsNormalizer"
}

func (n *AirconditionSpecsNormalizer) Filter(data map[string]string) map[string]string {

	for key, value := range data {

		if key == "Ισχύς Ψύξης" || key == "Ισχύς Θέρμανσης" {
			value = strings.ReplaceAll(value, " BTU", "")
		}

		if key == "Κατανάλωση Ψύξης" || key == "Κατανάλωση Θέρμανσης" {
			value = strings.ReplaceAll(value, " kWh/y", "")
		}

		if key == "Βαθμός Απόδοσης Ψύξης (SEER)" || key == "Βαθμός Απόδοσης Θέρμανσης (SCOP)" {
			value = strings.ReplaceAll(value, " W/W", "")
			value = strings.ReplaceAll(value, ",", ".")
		}
		data[key] = value
	}

	return data
}
