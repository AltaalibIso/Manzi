package controllers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"manzi/config"
	"net/http"
)

// RegisterHandler for /register
func RegisterHandler(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&user); err != nil || user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	collection := config.Client.Database(config.Config.Database.Name).Collection("users")

	err := collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		_, err := collection.InsertOne(context.TODO(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
		log.Printf("Error checking user existence: %v", err)
		return
	} else {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
	}
}
