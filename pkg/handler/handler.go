package handler

import (
	"errors"
	"fmt"
	"labs/pkg/command"
	"labs/pkg/models"
	"labs/pkg/worker"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	prodList []models.Product
	FilePath string

	commandList     []command.Command
	CommandFilePath string
)

func Init() {
	prodList = worker.CreateProductList(FilePath)
	commandList = command.LoadCommandsFromFile(CommandFilePath)
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

	if err := worker.WriteFile(FilePath, worker.CreateProductString(prodList)); err != nil {
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

	if err := worker.WriteFile(FilePath, worker.CreateProductString(prodList)); err != nil {
		log.Printf("Ошибка записи файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить данные"})
		return
	}

	c.Status(http.StatusNoContent)
}

func GetCommands(c *gin.Context) {
	c.JSON(http.StatusOK, commandList)
}

func AddCommand(c *gin.Context) {
	var input struct {
		Command string `json:"command"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd, err := command.ParseCommand(input.Command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	commandList = append(commandList, cmd)

	switch cmd.CommType {
	case "ADD":
		err = command.ExecuteAdd(cmd.Args, &prodList, FilePath)
	case "DEL":
		err = command.ExecuteDel(cmd.Args, &prodList, FilePath)
	case "SAVE":
		err = command.ExecuteSave(cmd.Args, &prodList)
	default:
		err = errors.New("неподдерживаемая команда")
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := command.SaveCommandsToFile(CommandFilePath, commandList); err != nil {
		log.Printf("Ошибка записи файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить данные"})
		return
	}

	c.JSON(http.StatusCreated, cmd)
}

func ExecuteAllCommands(c *gin.Context) {
	commands := make([]command.Command, len(commandList))
	copy(commands, commandList)

	for _, rec := range commands {
		var err error
		switch rec.CommType {
		case "ADD":
			err = command.ExecuteAdd(rec.Args, &prodList, FilePath)
		case "DEL":
			err = command.ExecuteDel(rec.Args, &prodList, FilePath)
		case "SAVE":
			err = command.ExecuteSave(rec.Args, &prodList)
		default:
			err = errors.New("неподдерживаемая команда")
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("ошибка выполнения команды %q: %v", rec.Args, err)})
			return
		}
	}

	c.JSON(http.StatusOK, prodList)
}
