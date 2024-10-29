package middleware

import (
	"github.com/gin-gonic/gin"
)

func CheckAuthServiceMiddleware(internalAuthToken string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.GetHeader("Internal_auth_token")
		if authToken != internalAuthToken {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		} else {
			ctx.Next()
		}
	}
}
