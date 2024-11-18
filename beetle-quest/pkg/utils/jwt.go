package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	jwt.StandardClaims
	Scope string `json:"scope"`
}

func (c CustomClaims) Valid() error {
	err := c.StandardClaims.Valid()
	if err != nil {
		return err
	}

	if c.Scope != "user" && c.Scope != "admin" {
		return fmt.Errorf("invalid scope")
	}

	return nil
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

func GenerateJWTToken(userID, scope string, secretKey []byte) (*jwt.Token, string, error) {
	now := time.Now()
	var expiresAt time.Time = time.Now()
	if scope == "user" {
		expiresAt = now.Add(time.Hour * 24)
	} else if scope == "admin" {
		expiresAt = now.Add(time.Hour)
	}

	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  "beetle-quest",
			Subject:   userID,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Issuer:    "beetle-quest",
			ExpiresAt: expiresAt.Unix(),
		},
		Scope: scope,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, "", err
	}
	return token, tokenString, nil
}
