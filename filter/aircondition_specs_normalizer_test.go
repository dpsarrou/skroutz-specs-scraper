package filter

import (
	"testing"
)

func TestItFiltersOutValueUnits(t *testing.T) {
	n := &AirconditionSpecsNormalizer{}
	data := map[string]string{
		"Ισχύς Ψύξης":                      "22000 BTU",
		"Ισχύς Θέρμανσης":                  "22000 BTU",
		"Κατανάλωση Ψύξης":                 "5 kWh/y",
		"Κατανάλωση Θέρμανσης":             "5 kWh/y",
		"Βαθμός Απόδοσης Ψύξης (SEER)":     "5 W/W",
		"Βαθμός Απόδοσης Θέρμανσης (SCOP)": "5,5 W/W",
		"This should not be altered":       "5,5 not altered",
	}

	filteredData := n.Filter(data)

	if filteredData["Ισχύς Ψύξης"] != "22000" {
		t.Fatalf("Could not filter out value")
	}
	if filteredData["Ισχύς Θέρμανσης"] != "22000" {
		t.Fatalf("Could not filter out value")
	}
	if filteredData["Κατανάλωση Ψύξης"] != "5" {
		t.Fatalf("Could not filter out value")
	}
	if filteredData["Κατανάλωση Θέρμανσης"] != "5" {
		t.Fatalf("Could not filter out value")
	}
	if filteredData["Βαθμός Απόδοσης Ψύξης (SEER)"] != "5" {
		t.Fatalf("Could not filter out value")
	}
	if filteredData["Βαθμός Απόδοσης Θέρμανσης (SCOP)"] != "5.5" {
		t.Fatalf("Could not filter out value")
	}
	if filteredData["This should not be altered"] != "5,5 not altered" {
		t.Fatalf("Value should not have been filtered, but it was")
	}
}
