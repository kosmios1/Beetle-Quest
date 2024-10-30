package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUserIDCorrespondWithSessionID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("user_id")
		checkUserID := ctx.GetHeader("user_id")
		if userID != checkUserID {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		} else {
			ctx.Next()
		}
	}
}
