package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUserID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenUserID, exists := ctx.Get("user_id")
		if tokenUserID == "" || !exists {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		paramUserID := ctx.Param("user_id")
		if paramUserID == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		if tokenUserID != paramUserID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		ctx.Next()
	}
}
