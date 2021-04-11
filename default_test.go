package vies_test

import (
	"testing"

	"github.com/macabu/vies"
)

func TestDefaultCheckVAT(t *testing.T) {
	vat, err := vies.CheckVAT("NL810060255B01")
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err.Error())
	}
	if vat == nil {
		t.Error("expected vat to not be nil, got nil")
	}
	if vat != nil && !vat.Valid {
		t.Error("expected vat to be true, got false")
	}
}

func TestDefaultCheckVATWithCountry(t *testing.T) {
	vat, err := vies.CheckVATWithCountry("NL", "NL810060255B01")
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err.Error())
	}
	if vat == nil {
		t.Error("expected vat to not be nil, got nil")
	}
	if vat != nil && !vat.Valid {
		t.Error("expected vat to be true, got false")
	}

	vat, err = vies.CheckVATWithCountry("NL", "810060255B01")
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err.Error())
	}
	if vat == nil {
		t.Error("expected vat to not be nil, got nil")
	}
	if vat != nil && !vat.Valid {
		t.Error("expected vat to be true, got false")
	}
}
