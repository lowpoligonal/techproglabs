package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func NewProduct(name string, count int, date time.Time) Product {
	return Product{
		ID:    uuid.New().String(),
		Name:  name,
		Count: count,
		Date:  date,
	}
}

func (p Product) ConvertToString() string {
	return fmt.Sprintf("%s %q %d %s", p.ID, p.Name, p.Count, p.Date.Format("2006-01-02"))
}

func ProductFromString(s string) (Product, error) {
	var dateStr string
	var p Product

	_, err := fmt.Sscanf(s, "%s %q %d %s", &p.ID, &p.Name, &p.Count, &dateStr)
	if err != nil {
		return p, fmt.Errorf("parse error: %w", err)
	}
	p.Date, _ = time.Parse("2006-01-02", dateStr)

	return p, nil
}
