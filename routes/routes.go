package routes

import (
	controller "github.com/ThembinkosiThemba/golang-crud/controllers"
	"github.com/gin-gonic/gin"
)

func BookRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/books", controller.GetBooks())
	incomingRoutes.GET("/books/:book_id", controller.GetBook())
	incomingRoutes.POST("/books/", controller.CreateBook())
	// incomingRoutes.PATCH("/books/:book_id", controller.UpdateBook())
}