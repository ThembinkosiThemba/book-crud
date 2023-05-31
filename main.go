package main

import (
	"os"

	"github.com/ThembinkosiThemba/golang-crud/database"
	"github.com/ThembinkosiThemba/golang-crud/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/joho/godotenv"
)

var bookCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()

	router.Use(gin.Logger())
	routes.BookRoutes(router)

	router.Run(":" + port)

}