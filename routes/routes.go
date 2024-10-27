package routes

import (
	"manzi/controllers"

	"github.com/gin-gonic/gin"
)

// InitUserRoutes initializes all user routes.
func InitUserRoutes(router *gin.Engine) {
	router.POST("/registration", controllers.RegisterHandler) // User registration route
}
