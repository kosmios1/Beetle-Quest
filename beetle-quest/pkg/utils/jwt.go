package utils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
}

func VerifyJWTToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
