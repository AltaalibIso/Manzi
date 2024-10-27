package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"manzi/repositories"
	"manzi/services"
	"net/http"
)

// RegisterHandler handles user registration.
func RegisterHandler(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Bind JSON input and validate required fields.
	if err := c.BindJSON(&user); err != nil || user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Validate username and password.
	if err := services.ValidateUsernamePassword(user.Username, user.Password); err != nil {
		// Check error type and return corresponding HTTP status.
		switch {
		case errors.Is(err, services.ErrInvalidUsername):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username. Must be 5-30 characters, with allowed characters: a-z, A-Z, 0-9, ., _, -"})
		case errors.Is(err, services.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password. Must be 5-30 characters, with allowed characters: a-z, A-Z, 0-9, ., _, -"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation error (username & password)"})
		}
		return
	}

	// Register the user.
	err := services.RegisterUser(user.Username, user.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		case errors.Is(err, repositories.ErrInvalidInput):
			c.JSON(http.StatusConflict, gin.H{"error": "Validation error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}
