package main

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/service"
	"beetle-quest/pkg/utils"
	"log"

	internalRepo "beetle-quest/internal/auth/repository"
	arepo "beetle-quest/pkg/repositories/serviceHttp/admin"
	urepo "beetle-quest/pkg/repositories/serviceHttp/user"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.GenOwnCertAndKey("auth")

	// This will connect to redis and return a store object used by the session middleware to store session data
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
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		AllowedHosts:          []string{},
	}))

	r.LoadHTMLGlob("templates/*")

	cnt := controller.NewAuthController(
		service.NewAuthService(urepo.NewUserRepo(), internalRepo.NewOauth2Repo(), arepo.NewAdminRepo()),
	)

	basePath := r.Group("/api/v1/auth")
	{
		basePath.POST("/register", cnt.Register)
		basePath.POST("/login", cnt.Login)
		basePath.GET("/logout", cnt.Logout)
		basePath.GET("/oauth2", cnt.Oauth2Callback)

		basePath.GET("/check_session", cnt.CheckSession)
		basePath.Any("/traefik/verify", cnt.Verify)
	}

	adminSpecific := r.Group("/api/v1/auth/admin")
	{
		adminSpecific.POST("/login", cnt.AdminLogin)
		adminSpecific.GET("/oauth2", cnt.AdminOauth2Callback)
	}

	server := utils.SetupHTPPSServer(r)
	if err := server.ListenAndServeTLS("/serverCert.pem", "/serverKey.pem"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
