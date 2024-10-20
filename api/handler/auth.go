package handler

import (
	"cms-api/api/models"
	"cms-api/internal/database/connection"
	"cms-api/pkg/utility/common"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

// LoginRequest is the request structure for login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates the user
func Login(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		common.GenerateErrorResponse(c, "Invalid input")
		return
	}

	// Pass the context to isValidUser and get user info
	user, isValid := isValidUser(c, loginReq.Username, loginReq.Password)
	if !isValid {
		return // Error response already generated in isValidUser
	}

	// Generate tokens if user is valid
	accessToken, err := common.GenerateAccessToken(user.Username, user.Role) // Use user's role
	if err != nil {
		common.GenerateErrorResponse(c, "Could not create access token")
		return
	}
	refreshToken, err := common.GenerateRefreshToken(user.Username, user.Role) // Use user's role
	if err != nil {
		common.GenerateErrorResponse(c, "Could not create refresh token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

// RefreshToken updates the access token
func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GenerateErrorResponse(c, "Invalid input")
		return
	}

	claims := &common.RefreshClaims{}
	tkn, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(viper.GetString("jwt_secret_key")), nil
	})

	// Check if the token is valid and the issuer is correct
	if err != nil || !tkn.Valid || claims.Issuer != "refresh" {
		common.GenerateErrorResponse(c, "invalid or expired refresh token")
		return
	}

	// Create a new access token
	accessToken, err := common.GenerateAccessToken(claims.Username, claims.Role)
	if err != nil {
		common.GenerateErrorResponse(c, "Could not create access token")
		return
	}
	refreshToken, err := common.GenerateRefreshToken(claims.Username, claims.Role)
	if err != nil {
		common.GenerateErrorResponse(c, "Could not create refresh token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}
func isValidUser(ctx *gin.Context, username, password string) (models.User, bool) {
	var user models.User

	// Fetch the user from the database
	if err := connection.DB.Where("username = ?", username).First(&user).Error; err != nil {
		common.GenerateErrorResponse(ctx, "Invalid username")
		return models.User{}, false
	}

	// Check if the provided password matches the hashed password stored in the database
	if err := common.CheckPasswordHash(password, user.Password); err != nil {
		common.GenerateErrorResponse(ctx, "Invalid password")
		return models.User{}, false
	}

	// User validation successful
	return user, true
}
