package routes

import (
	"booklist-back/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoutes(router *gin.Engine) {
	router.POST("/books", controllers.CreateBook)
	router.GET("/books", controllers.ListBooks)
	router.PUT("/books/:id", controllers.UpdateBook)
	router.DELETE("/books/:id", controllers.DeleteBook)

}
