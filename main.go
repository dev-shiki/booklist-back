package main

import (
	"booklist-back/database"
	"booklist-back/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()
	defer database.DB.Close()

	router := gin.Default()

	router.Use(cors.Default())

	routes.BookRoutes(router)
	routes.CategoryRoutes(router)

	router.Run(":8080")
}
