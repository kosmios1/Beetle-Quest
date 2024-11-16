package middleware

import (
	"beetle-quest/pkg/utils"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	jwtSecretKey = utils.PanicIfError[[]byte](hex.DecodeString(utils.FindEnv("JWT_SECRET_KEY")))
)

func CheckAdminJWTAuthorizationToken() gin.HandlerFunc {
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

		if claims.Valid() != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims["scope"] != "admin" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("admin_id", claims["sub"])
		ctx.Next()
	}
}
