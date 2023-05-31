package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID        	primitive.ObjectID    `bson:"_id,omitempty"`
	Title     	string    			`json:"title,omitempty"`
	Genre     	string    			`json:"genre,omitempty"`
	Year      	string    			`json:"year,omitempty"`
	Created_At 	time.Time 			`json:"created_at,omitempty"`
	Updated_At 	time.Time 			`json:"updated_at,omitempty"`
	Book_id     string            `json:"book_id"`
}