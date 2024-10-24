package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"manzi/services"
	"net/http"
)

func RegisterHandler(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&user); err != nil || user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	if err := services.ValidateUsernamePassword(user.Username, user.Password); err != nil {
		// Проверяем тип ошибки и возвращаем соответствующий HTTP-статус
		switch {
		case errors.Is(err, services.ErrInvalidUsername):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username. Must be 5-30 characters, with allowed characters: a-z, A-Z, 0-9, ., _, -"})
		case errors.Is(err, services.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password. Must be 5-30 characters, with allowed characters: a-z, A-Z, 0-9, ., _, -"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation error"})
		}
		return
	}

	err := services.RegisterUser(user.Username, user.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}
