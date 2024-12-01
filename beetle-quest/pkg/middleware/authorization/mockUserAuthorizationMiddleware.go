//go:build beetleQuestTest

package middleware

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	jwtSecretKey = utils.PanicIfError[[]byte](hex.DecodeString(utils.FindEnv("JWT_SECRET_KEY")))
)

func CheckJWTAuthorizationToken(requestedScope models.Scope) gin.HandlerFunc {
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
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})

		if err != nil {
			ctx.HTML(http.StatusUnauthorized, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims["sub"])
		ctx.Next()
	}
}
