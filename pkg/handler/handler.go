package handler

import (
	"labs/pkg/models"
	"labs/pkg/worker"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const filePath = "D:/ucheba/techproglabs/input.txt"

var prodList []models.Product

func Init() {
	prodList = worker.CreateProductList(filePath)
}

func GetProducts(c *gin.Context) {
	c.JSON(http.StatusOK, prodList)
}

func AddProduct(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		Count    int    `json:"count"`
		Date     string `json:"date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Count < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Количество не может быть отрицательным"})
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты"})
		return
	}

	newProd := models.NewProduct(input.Name, input.Category, input.Count, date)
	prodList = append(prodList, newProd)

	if err := worker.WriteFile(filePath, worker.CreateProductString(prodList)); err != nil {
		log.Printf("Ошибка записи файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить данные"})
		return
	}

	c.JSON(http.StatusCreated, newProd)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	index := -1
	for i, p := range prodList {
		if p.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	prodList = append(prodList[:index], prodList[index+1:]...)

	if err := worker.WriteFile(filePath, worker.CreateProductString(prodList)); err != nil {
		log.Printf("Ошибка записи файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить данные"})
		return
	}

	c.Status(http.StatusNoContent)
}
