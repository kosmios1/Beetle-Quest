//go:build beetleQuestTest

package middleware

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func NewSecureMiddleware() gin.HandlerFunc {
	return secure.New(secure.Config{
		IsDevelopment: true,
	})
}
