package main

import (
	"fmt"
	"time"
)

type Product struct {
	Name  string
	Count int
	Date  time.Time
}

func main() {
	var product Product
	var inputDate string

	fmt.Println("name, count, date")
	n, err := fmt.Scanf(`%q %d %v`, &product.Name, &product.Count, &inputDate)
	product.Date, _ = time.Parse("2006-01-02", inputDate)

	fmt.Printf("Name: %s\nCount: %d\nDate: %v\n", product.Name, product.Count, product.Date)

	fmt.Println(n, err)
}
