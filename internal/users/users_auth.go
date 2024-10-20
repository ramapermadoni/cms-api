package users

import (
	"cms-api/api/handler"

	"github.com/gin-gonic/gin"
)

func AuthInitiator(router *gin.Engine) {
	router.POST("/login", handler.Login)                // Login
	router.POST("/refresh-token", handler.RefreshToken) // Refresh Token
	// Add any other routes necessary for the user module
}
