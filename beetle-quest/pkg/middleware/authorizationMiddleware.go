package middleware

import (
	"beetle-quest/pkg/utils"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	jwtSecretKey = utils.PanicIfError[[]byte](hex.DecodeString(os.Getenv("JWT_SECRET_KEY")))
)

func CheckJWTAuthorizationToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("access_token")
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := utils.VerifyJWTToken(cookie.Value, jwtSecretKey)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user_id", claims["sub"])
		ctx.Next()
	}
}
