package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAuthServiceMiddleware(internalAuthToken string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.GetHeader("Internal_auth_token")
		if authToken != internalAuthToken {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Redirect(http.StatusMovedPermanently, "/api/v1/login")
			ctx.Abort()
			return
		} else {
			ctx.Next()
		}
	}
}
