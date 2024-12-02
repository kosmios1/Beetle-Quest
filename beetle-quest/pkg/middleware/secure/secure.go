//go:build !beetleQuestTest

package middleware

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func NewSecureMiddleware() gin.HandlerFunc {
	return secure.New(secure.Config{
		SSLRedirect:          true,
		IsDevelopment:        false,
		STSSeconds:           315360000,
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		// ContentSecurityPolicy: "default-src 'self'; img-src 'self'; script-src 'self' 'unsafe-eval'; style-src 'self'",
		IENoOpen:        true,
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		AllowedHosts:    []string{},
	})
}
