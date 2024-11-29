package main

import (
	"beetle-quest/pkg/utils"
	"log"

	entrypoint "beetle-quest/internal/auth/entrypoints"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.GenOwnCertAndKey("auth-service")

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(secure.New(secure.Config{
		SSLRedirect:           true,
		IsDevelopment:         false,
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; img-src 'self'; script-src 'self' 'unsafe-eval'; style-src 'self'",
		IENoOpen:              true,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		AllowedHosts:          []string{},
	}))
	r.LoadHTMLGlob("templates/*")

	cnt := entrypoint.NewAuthController()
	basePath := r.Group("/api/v1/auth")
	{
		basePath.GET("/authPage", cnt.AuthenticationPage)
		basePath.GET("/tokenPage", cnt.TokenPage)
		basePath.GET("/authorizePage", cnt.AuthorizePage)

		basePath.POST("/register", cnt.Register)
		basePath.POST("/login", cnt.Login)
		basePath.GET("/logout", cnt.Logout)

		basePath.GET("/check_session", cnt.CheckSession)
		basePath.Any("/traefik/verify", cnt.Verify)

		oauthPath := basePath.Group("/oauth")
		{
			oauthPath.POST("/authorize", cnt.OauthAuthorize)
			oauthPath.POST("/token", cnt.OauthToken)
		}
	}

	adminSpecific := r.Group("/api/v1/auth/admin")
	{
		adminSpecific.POST("/login", cnt.AdminLogin)
	}

	server := utils.SetupHTPPSServer(r)
	if err := server.ListenAndServeTLS("/serverCert.pem", "/serverKey.pem"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
