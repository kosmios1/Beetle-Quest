package main

import (
	"beetle-quest/internal/user/controller"
	"beetle-quest/internal/user/repository"
	"beetle-quest/internal/user/service"
	"beetle-quest/pkg/middleware"
	"beetle-quest/pkg/utils"
	"log"

	internalMiddleware "beetle-quest/internal/user/middleware"
	gHttpGrepo "beetle-quest/pkg/repositories/serviceHttp/gacha"
	mHttpGrepo "beetle-quest/pkg/repositories/serviceHttp/market"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.GenOwnCertAndKey("user")

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

	cnt := controller.NewUserController(
		service.NewUserService(repository.NewUserRepo(), gHttpGrepo.NewGachaRepo(), mHttpGrepo.NewMarketRepo()),
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
