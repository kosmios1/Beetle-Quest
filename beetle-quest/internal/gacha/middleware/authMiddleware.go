package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		sessionID := session.Get("session_id")
		if sessionID == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Redirect(http.StatusMovedPermanently, "/api/v1/login")
			ctx.Abort()
			return
		} else {
			ctx.Next()
		}
	}
}
