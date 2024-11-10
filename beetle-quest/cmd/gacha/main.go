package main

import (
	"beetle-quest/internal/gacha/controller"
	"beetle-quest/internal/gacha/repository"
	"beetle-quest/internal/gacha/service"
	"beetle-quest/pkg/middleware"

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

	cnt := controller.NewGachaController(
		service.NewGachaService(repository.NewGachaRepo()),
	)

	basePath := r.Group("/api/v1/gacha")
	basePath.Use(middleware.CheckJWTAuthorizationToken())
	{
		basePath.POST("/roll", cnt.Roll)
		basePath.GET("/list", cnt.List)
		basePath.GET("/:gacha_id", cnt.GetGachaDetails)
	}

	r.Run()
}
