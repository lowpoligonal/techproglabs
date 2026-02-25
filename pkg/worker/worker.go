package worker

import (
	"fmt"
	"io"
	"labs/pkg/models"
	"os"
	"strings"
)

func ReadFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data := make([]byte, 64)
	fileData := ""

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		fileData += string(data[:n])
	}
	return fileData
}

func WriteFile(path, data string) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	file.WriteString(data)
}

func CreateProductList(filePath string) []models.Product {
	data := ReadFile(filePath)
	lines := strings.Split(data, "\n")
	var list []models.Product

	for _, line := range lines {
		p, err := models.ProductFromString(line)
		if err != nil {
		}
		list = append(list, p)
	}
	return list
}

func CreateProductString(data []models.Product) string {
	var output string
	for _, prod := range data {
		output += prod.ConvertToString()
		output += "\n"
	}
	return output
}
