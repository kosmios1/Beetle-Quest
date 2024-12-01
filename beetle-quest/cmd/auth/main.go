package main

import (
	"beetle-quest/pkg/httpserver"
	"log"

	entrypoint "beetle-quest/internal/auth/entrypoints"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	packageMiddleware "beetle-quest/pkg/middleware/secure"
)

func main() {
	httpserver.GenOwnCertAndKey("auth-service")

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(packageMiddleware.NewSecureMiddleware())

	r.Use(cors.Default())
	r.LoadHTMLGlob("templates/*")

	cnt, err := entrypoint.NewAuthController()
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}

	oauthPath := r.Group("/oauth")
	{
		oauthPath.GET("/authorize", cnt.OauthAuthorize)
		oauthPath.POST("/token", cnt.OauthToken)
	}

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
	}

	adminSpecific := r.Group("/api/v1/auth/admin")
	{
		adminSpecific.POST("/login", cnt.AdminLogin)
	}

	server := httpserver.SetupHTPPSServer(r)
	httpserver.ListenAndServe(server)
}
