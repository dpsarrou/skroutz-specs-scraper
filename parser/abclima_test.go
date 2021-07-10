package parser

import (
	"os"
	"testing"
)

func TestItParsesAbClimaData(t *testing.T) {
	d, err := os.ReadFile("../testdata/shop_3292_abclima.html")
	if err != nil {
		t.Fatal(err)
	}
	p := &AbClimaParser{}
	data, err := p.ParseProductSpecs(d)
	if err != nil {
		t.Error(err)
	}
	if data["Ελάχιστη Ψύξη (Btu/h)"] != "5802" {
		t.Errorf("Could not parse Ελάχιστη ψύξη (Btu/h)")
	}
	if data["Μέγιστη Ψύξη (Btu/h)"] != "20478" {
		t.Errorf("Could not parse Μέγιστη Ψύξη (Btu/h)")
	}
}
