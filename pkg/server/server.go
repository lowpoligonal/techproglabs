package server

import (
	"labs/pkg/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	handler.Init()
	gin.SetMode(gin.ReleaseMode)

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

		api.GET("/commands", handler.GetCommands)
		api.POST("/commands", handler.AddCommand)
		//api.DELETE("/commands/:id", handler.DeleteCommand)

		api.POST("/commands/execute-all", handler.ExecuteAllCommands)
	}

	if err := router.Run(":8088"); err != nil {
		log.Fatal("Server failed:", err)
	}
}
