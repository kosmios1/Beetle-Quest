package main

import (
	"beetle-quest/pkg/utils"
	"log"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"

	"beetle-quest/internal/admin/controller"
	"beetle-quest/internal/admin/middleware"
	"beetle-quest/internal/admin/repository"
	"beetle-quest/internal/admin/service"
)

func main() {
	utils.GenOwnCertAndKey("admin")

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

	cnt := controller.NewAdminController(
		service.NewAdminService(repository.NewAdminRepo()),
	)

	basePath := r.Group("/api/v1/admin")
	basePath.Use(middleware.CheckAdminJWTAuthorizationToken())
	{
		basePath.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"Message": "Admin API v1 is up and running!"})
		})
	}

	internalPath := r.Group("/api/v1/internal/admin")
	{
		internalPath.POST("/find_by_id", cnt.FindByID)
	}

	server := utils.SetupHTPPSServer(r)
	if err := server.ListenAndServeTLS("/serverCert.pem", "/serverKey.pem"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
