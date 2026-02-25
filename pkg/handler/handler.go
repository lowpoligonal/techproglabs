package handler

import (
	"labs/pkg/models"
	"labs/pkg/worker"
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
		Name  string `json:"name" binding:"required"`
		Count int    `json:"count" binding:"required"`
		Date  string `json:"date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, _ := time.Parse("2006-01-02", input.Date)

	newProd := models.NewProduct(input.Name, input.Count, date)
	prodList = append(prodList, newProd)

	worker.WriteFile(filePath, worker.CreateProductString(prodList))

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

	worker.WriteFile(filePath, worker.CreateProductString(prodList))

	c.Status(http.StatusNoContent)
}
