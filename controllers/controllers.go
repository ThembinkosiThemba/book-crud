package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ThembinkosiThemba/golang-crud/database"
	"github.com/ThembinkosiThemba/golang-crud/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var bookCollection *mongo.Collection = database.OpenCollection(database.Client, "book")
var validate = validator.New()

func GetBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var allBooks []bson.M

		result, err := bookCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing all book items"})
		}

		if err = result.All(ctx, &allBooks); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allBooks)

	}
}

func GetBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		bookId := c.Param("book_id")
		var book models.Book

		err := bookCollection.FindOne(ctx, bson.M{"book_id": bookId})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while getting book item"})
		}

		c.JSON(http.StatusOK, book)

	}
}

func CreateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var book models.Book

		// JSON unmashaling
		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validatorErr := validate.Struct(book)

		if validatorErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validatorErr.Error()})
			return
		}

		book.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		book.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		book.ID = primitive.NewObjectID()
		book.Book_id = book.ID.Hex()

		result, insertErr := bookCollection.InsertOne(ctx, book)
		if insertErr != nil {
			msg := fmt.Sprintf("Book item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
		defer cancel()
	}
}

func UpdateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var book models.Book

		bookId := c.Param("book_id")

		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		var updateObj primitive.D

		if book.Title != "" {
			updateObj = append(updateObj, bson.E{Key: "title", Value: book.Title})
		}

		if book.Genre != "" {
			updateObj = append(updateObj, bson.E{Key: "genre", Value: book.Genre})
		}

		if book.Year != "" {
			updateObj = append(updateObj, bson.E{Key: "year", Value: book.Year})
		}

		book.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: book.Updated_At})

		upsert := true
		filter := bson.M{"book_id": bookId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := bookCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := "Menu uodate failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}

		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}
