package routes

import (
	"manzi/controllers"

	"github.com/gin-gonic/gin"
)

// InitUserRoutes init all user routs
func InitUserRoutes(router *gin.Engine) {
	router.POST("/registration", controllers.RegisterHandler)
}
