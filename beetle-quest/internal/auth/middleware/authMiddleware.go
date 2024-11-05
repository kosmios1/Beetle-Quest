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

		// TODO: Check that the user ID is valid, otherwise an active session of a deleted user can be used to access the API

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
