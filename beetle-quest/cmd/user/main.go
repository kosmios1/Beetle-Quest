package main

import (
	"beetle-quest/pkg/middleware"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"log"

	entrypoint "beetle-quest/internal/user/entrypoints"
	internalMiddleware "beetle-quest/internal/user/middleware"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.GenOwnCertAndKey("user-service")

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

	cnt := entrypoint.NewUserController()

	basePath := r.Group("/api/v1/user")
	basePath.Use(middleware.CheckJWTAuthorizationToken(models.UserScope))
	{
		accountGroup := basePath.Group("/account")
		accountGroup.Use(internalMiddleware.CheckUserID())
		{
			accountGroup.GET("/:user_id", cnt.GetUserAccountDetails)
			accountGroup.PATCH("/:user_id", cnt.UpdateUserAccountDetails)
			accountGroup.POST("/:user_id", cnt.DeleteUserAccount)
		}
	}

	internalPath := r.Group("/api/v1/internal/user")
	{
		internalPath.GET("/get_all", cnt.GetAllUsers)
		internalPath.POST("/create", cnt.CreateUser)
		internalPath.POST("/update", cnt.UpdateUser)

		internalPath.POST("/find_by_id", cnt.FindByID)
		internalPath.POST("/find_by_username", cnt.FindByUsername)
	}

	server := utils.SetupHTPPSServer(r)
	if err := server.ListenAndServeTLS("/serverCert.pem", "/serverKey.pem"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
