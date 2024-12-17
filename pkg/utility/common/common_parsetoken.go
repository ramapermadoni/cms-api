package common

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

// ParseToken parses and validates JWT token, returns claims if valid
func ParseToken(tokenString string) (*AccessClaims, error) {
	claims := &AccessClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt_secret_key")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
