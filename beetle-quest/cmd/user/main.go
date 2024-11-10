package main

import (
	"beetle-quest/internal/user/controller"
	"beetle-quest/internal/user/repository"
	"beetle-quest/internal/user/service"
	"beetle-quest/pkg/middleware"

	internalMiddleware "beetle-quest/internal/user/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	// TODO: Uncomment this when having a valid SSL certificate
	// r.Use(secure.New(secure.Config{
	// 	SSLRedirect:           true,
	// 	IsDevelopment:         false,
	// 	STSSeconds:            315360000,
	// 	STSIncludeSubdomains:  true,
	// 	FrameDeny:             true,
	// 	ContentTypeNosniff:    true,
	// 	BrowserXssFilter:      true,
	// 	ContentSecurityPolicy: "default-src 'self'",
	// 	IENoOpen:              true,
	// 	SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	// 	AllowedHosts:          []string{},
	// }))

	r.LoadHTMLGlob("templates/*")

	cnt := controller.NewUserController(
		service.NewUserService(repository.NewUserRepo()),
	)

	basePath := r.Group("/api/v1/user")
	basePath.Use(middleware.CheckJWTAuthorizationToken())
	{
		accountGroup := basePath.Group("/account")
		accountGroup.Use(internalMiddleware.CheckUserID())
		{
			accountGroup.GET("/:user_id", cnt.GetUserAccountDetails)
			accountGroup.PATCH("/:user_id", cnt.UpdateUserAccountDetails)
			accountGroup.DELETE("/:user_id", cnt.DeleteUserAccount)
		}

		basePath.GET("/:user_id/gacha", cnt.GetUserGachaList)
		basePath.GET("/:user_id/gacha/:gacha_id", cnt.GetUserGachaDetails)
	}

	internalPath := r.Group("/api/v1/internal/user")
	// TODO: This can be reached only within the microservices network
	{
		internalPath.POST("/create", cnt.CreateUser)
		internalPath.POST("/find_by_id", cnt.FindByID)
		internalPath.POST("/find_by_username", cnt.FindByUsername)
	}

	r.Run()
}
