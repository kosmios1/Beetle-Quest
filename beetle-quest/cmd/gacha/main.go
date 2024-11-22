package main

import (
	"beetle-quest/internal/gacha/controller"
	"beetle-quest/internal/gacha/repository"
	"beetle-quest/internal/gacha/service"
	"beetle-quest/pkg/middleware"
	"beetle-quest/pkg/utils"
	"log"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.GenOwnCertAndKey("gacha-service")

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

	cnt := controller.NewGachaController(
		service.NewGachaService(repository.NewGachaRepo()),
	)

	basePath := r.Group("/api/v1/gacha")
	basePath.Use(middleware.CheckJWTAuthorizationToken())
	{
		basePath.GET("/list", cnt.List)
		basePath.GET("/:gacha_id", cnt.GetGachaDetails)

		basePath.GET("/user/:user_id/list", cnt.GetUserGachaList)
		basePath.GET("/:gacha_id/user/:user_id", cnt.GetUserGachaDetails)
	}

	internalPath := r.Group("/api/v1/internal/gacha")
	{
		internalPath.POST("/create", cnt.CreateGacha)
		internalPath.POST("/update", cnt.UpdateGacha)
		internalPath.POST("/delete", cnt.DeleteGacha)

		internalPath.POST("/get_user_gachas", cnt.GetUserGachas)
		internalPath.POST("/remove_user_gachas", cnt.RemoveUserGachas)
		internalPath.POST("/add_gacha_to_user", cnt.AddGachaToUser)
		internalPath.POST("/remove_gacha_from_user", cnt.RemoveGachaFromUser)
		internalPath.POST("/find_by_id", cnt.FindByID)
		internalPath.GET("/get_all", cnt.GetAll)
	}

	server := utils.SetupHTPPSServer(r)
	if err := server.ListenAndServeTLS("/serverCert.pem", "/serverKey.pem"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
