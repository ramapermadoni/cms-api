package middlewares

import (
	"cms-api/pkg/utility/common"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

// AccessClaims for access token
// type AccessClaims struct {
// 	Issuer   string `json:"iss"`
// 	ID       int    `json:"id"`
// 	Username string `json:"username"`
// 	Role     string `json:"role"` // Tambahkan field Role
// 	jwt.RegisteredClaims
// }

// JwtMiddleware is a middleware for JWT authentication
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := GetJwtTokenFromHeader(c)
		if err != nil {
			common.GenerateErrorResponse(c, err.Error())
			c.Abort()
			return
		}

		claims := &common.AccessClaims{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(viper.GetString("jwt_secret_key")), nil
		})

		if err != nil || !tkn.Valid {
			log.Printf("Invalid token: %v", err)
			common.GenerateErrorResponse(c, "invalid or expired access token")
			c.Abort()
			return
		}

		// Check if token is an access token (optional: you can add claims to identify the type of token)
		if !isAccessToken(tokenString) {
			common.GenerateErrorResponse(c, "token is not allowed for this route")
			c.Abort()
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			common.GenerateErrorResponse(c, "access token expired")
			c.Abort()
			return
		}

		c.Set("id_user", claims.ID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role) // Simpan role ke context
		log.Printf("User id_user set in context: %d", claims.ID)
		log.Printf("User username set in context: %s", claims.Username)
		log.Printf("User role set in context: %s", claims.Role)
		c.Next()
	}
}

// GetJwtTokenFromHeader retrieves the JWT token from the Authorization header
func GetJwtTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if common.IsEmptyField(authHeader) {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	return parts[1], nil
}

// isAccessToken checks if the provided token is an access token
func isAccessToken(token string) bool {
	claims := &common.AccessClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt_secret_key")), nil
	})

	if err != nil {
		return false
	}

	// Check if the issuer is valid
	if claims.Issuer != "access" { // Replace with your actual issuer
		return false
	}

	return claims.RegisteredClaims.ExpiresAt.Time.After(time.Now())
}

// RoleMiddleware checks the user's role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		log.Printf("allowedRoles: %v", allowedRoles)
		log.Printf("role: %v", role)

		if !exists {
			common.GenerateErrorResponse(c, "role not found in context")
			c.Abort()
			return
		}

		// Lakukan type assertion
		roleStr, ok := role.(string)
		if !ok {
			common.GenerateErrorResponse(c, "invalid role type")
			c.Abort()
			return
		}

		// Cek apakah role pengguna ada dalam allowedRoles
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		common.GenerateErrorResponse(c, "access denied")
		c.Abort()
	}
}
