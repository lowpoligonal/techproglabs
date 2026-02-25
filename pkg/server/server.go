package server

import (
	"labs/pkg/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	handler.Init()

	router := gin.Default()

	router.Static("/static", "./web/static")

	router.GET("/", func(c *gin.Context) {
		c.File("./web/static/index.html")
	})

	api := router.Group("/api")
	{
		api.GET("/products", handler.GetProducts)
		api.POST("/products", handler.AddProduct)
		api.DELETE("/products/:id", handler.DeleteProduct)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed:", err)
	}
}
