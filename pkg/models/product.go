package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
	Count    int       `json:"count"`
}

func NewProduct(name, category string, count int, date time.Time) Product {
	return Product{
		ID:       uuid.New().String(),
		Name:     name,
		Count:    count,
		Date:     date,
		Category: category,
	}
}

func (p Product) ConvertToString() string {
	return fmt.Sprintf("%s %q %s %d %s", p.ID, p.Name, p.Category, p.Count, p.Date.Format("2006-01-02"))
}

func ProductFromString(s string) (Product, error) {
	var dateStr string
	var p Product

	_, err := fmt.Sscanf(s, "%s %q %s %d %s", &p.ID, &p.Name, &p.Category, &p.Count, &dateStr)
	if err != nil {
		return p, fmt.Errorf("Некорректная строка: %w", err)
	}

	p.Date, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		return p, fmt.Errorf("Ошибка парсинга даты: %w", err)
	}

	return p, nil
}
