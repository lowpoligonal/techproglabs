package models

import (
	"testing"
	"time"
)

func TestNewProduct(t *testing.T) {
	date := time.Date(2025, time.March, 11, 0, 0, 0, 0, time.UTC)
	p := NewProduct("Milk", "Dairy", 10, date)

	if p.ID == "" {
		t.Error("expected non-empty ID")
	}
	if p.Name != "Milk" {
		t.Errorf("expected name Milk, got %s", p.Name)
	}
	if p.Category != "Dairy" {
		t.Errorf("expected category Dairy, got %s", p.Category)
	}
	if p.Count != 10 {
		t.Errorf("expected count 10, got %d", p.Count)
	}
	if !p.Date.Equal(date) {
		t.Errorf("expected date %v, got %v", date, p.Date)
	}
}

func TestProduct_ConvertToString(t *testing.T) {
	date := time.Date(2025, time.March, 11, 0, 0, 0, 0, time.UTC)
	p := Product{
		ID:       "test-id",
		Name:     "Bread",
		Category: "Bakery",
		Count:    3,
		Date:     date,
	}
	expected := `test-id "Bread" Bakery 3 2025-03-11`
	result := p.ConvertToString()
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestProductFromString_Valid(t *testing.T) {
	line := `123e4567-e89b-12d3-a456-426614174000 "Milk" Dairy 10 2025-03-11`
	p, err := ProductFromString(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("wrong ID: %s", p.ID)
	}
	if p.Name != "Milk" {
		t.Errorf("wrong Name: %s", p.Name)
	}
	if p.Category != "Dairy" {
		t.Errorf("wrong Category: %s", p.Category)
	}
	if p.Count != 10 {
		t.Errorf("wrong Count: %d", p.Count)
	}
	expectedDate, _ := time.Parse("2006-01-02", "2025-03-11")
	if !p.Date.Equal(expectedDate) {
		t.Errorf("wrong Date: %v", p.Date)
	}
}

func TestProductFromString_InvalidFormat(t *testing.T) {
	line := `invalid line without quotes`
	_, err := ProductFromString(line)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestProductFromString_InvalidDate(t *testing.T) {
	line := `123 "Milk" Dairy 10 2025/03/11`
	_, err := ProductFromString(line)
	if err == nil {
		t.Error("expected error for invalid date, got nil")
	}
}

func TestProductFromString_MissingField(t *testing.T) {
	line := `123 "Milk" Dairy 10` // пропущена дата
	_, err := ProductFromString(line)
	if err == nil {
		t.Error("expected error for missing field, got nil")
	}
}

func TestProductFromString_ExtraSpaces(t *testing.T) {
	line := `  123  "Milk"   Dairy   10   2025-03-11  `
	p, err := ProductFromString(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != "123" {
		t.Errorf("expected ID 123, got %q", p.ID)
	}
}
