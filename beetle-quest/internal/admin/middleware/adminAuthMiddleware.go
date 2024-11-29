package middleware

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	jwtSecretKey = utils.PanicIfError[[]byte](hex.DecodeString(utils.FindEnv("JWT_SECRET_KEY")))
)

func CheckAdminJWTAuthorizationToken(requestedScope models.Scope) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.HTML(http.StatusUnauthorized, "errorMsg.tmpl", gin.H{"Error": "Unauthorized access."})
			ctx.Abort()
			return
		}

		parsedAuthHeader := strings.Split(authHeader, " ")
		if len(parsedAuthHeader) != 2 {
			ctx.HTML(http.StatusUnauthorized, "errorMsg.tmpl", gin.H{"Error": "Unauthorized access."})
			ctx.Abort()
			return
		}

		accessToken := parsedAuthHeader[1]
		claims, err := utils.VerifyJWTToken(accessToken, jwtSecretKey)
		if err != nil {
			ctx.HTML(http.StatusUnauthorized, "errorMsg.tmpl", gin.H{"Error": "Unauthorized access."})
			ctx.Abort()
			return
		}

		if claims.Valid() != nil {
			ctx.HTML(http.StatusUnauthorized, "errorMsg.tmpl", gin.H{"Error": "Unauthorized access."})
			ctx.Abort()
			return
		}

		scope := claims["scope"].(string)
		scopes := strings.Split(scope, ", ")
		for _, s := range scopes {
			if models.Scope(s) == models.AdminScope {
				ctx.Set("admin_id", claims["sub"])
				ctx.Next()
				return
			}
		}

		ctx.HTML(http.StatusUnauthorized, "errorMsg.tmpl", gin.H{"Error": "Unauthorized access."})
		ctx.Abort()
	}
}
