package worker

import (
	"io"
	"labs/pkg/models"
	"log"
	"os"
	"strings"
)

func ReadFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data := make([]byte, 64)
	fileData := ""

	for {
		n, err := file.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		fileData += string(data[:n])
	}
	return fileData, nil
}

func WriteFile(path, data string) error {
	return os.WriteFile(path, []byte(data), 0644)
}

func CreateProductList(filePath string) []models.Product {
	data, err := ReadFile(filePath)
	if err != nil {
		log.Printf("Ошибка чтения файла %s: %v", filePath, err)
		return nil
	}

	lines := strings.Split(data, "\n")

	list := make([]models.Product, 0)

	if len(lines) == 0 {
		return list
	}

	for _, line := range lines {
		p, err := models.ProductFromString(line)
		if err != nil {
			log.Printf("Ошибка: %v", err)
			continue
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
