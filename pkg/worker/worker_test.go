package worker

import (
	"labs/pkg/models"
	"os"
	"testing"
	"time"
)

func TestReadFile(t *testing.T) {
	content := "line1\nline2\n"
	tmpfile, err := os.CreateTemp("", "testread*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	data, err := ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != content {
		t.Errorf("expected %q, got %q", content, data)
	}
}

func TestReadFile_NotExist(t *testing.T) {
	_, err := ReadFile("nonexistent_file.txt")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestWriteFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testwrite*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	data := "test data"
	err = WriteFile(tmpfile.Name(), data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != data {
		t.Errorf("expected %q, got %q", data, string(content))
	}
}

func TestCreateProductList(t *testing.T) {
	content := `id1 "Product1" Cat1 10 2025-03-11
this is invalid line
id2 "Product2" Cat2 5 2025-03-12`
	tmpfile, err := os.CreateTemp("", "testlist*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	list := CreateProductList(tmpfile.Name())
	if len(list) != 2 {
		t.Errorf("expected 2 products, got %d", len(list))
	}
	if len(list) > 0 && list[0].ID != "id1" {
		t.Errorf("expected first product ID id1, got %q", list[0].ID)
	}
	if len(list) > 1 && list[1].ID != "id2" {
		t.Errorf("expected second product ID id2, got %q", list[1].ID)
	}
}

func TestCreateProductList_EmptyFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testempty*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	list := CreateProductList(tmpfile.Name())
	if list == nil {
		t.Error("expected empty slice, got nil")
	}
	if len(list) != 0 {
		t.Errorf("expected 0 products, got %d", len(list))
	}
}

func TestCreateProductList_FileNotExist(t *testing.T) {
	list := CreateProductList("nonexistent_file.txt")
	if list != nil {
		t.Errorf("expected nil, got %v", list)
	}
}

func TestCreateProductString(t *testing.T) {
	date1, _ := time.Parse("2006-01-02", "2025-03-11")
	date2, _ := time.Parse("2006-01-02", "2025-03-12")
	products := []models.Product{
		{ID: "id1", Name: "Prod1", Category: "Cat1", Count: 10, Date: date1},
		{ID: "id2", Name: "Prod2", Category: "Cat2", Count: 5, Date: date2},
	}
	expected := `id1 "Prod1" Cat1 10 2025-03-11
id2 "Prod2" Cat2 5 2025-03-12
`
	result := CreateProductString(products)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestCreateProductString_Empty(t *testing.T) {
	result := CreateProductString([]models.Product{})
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}
