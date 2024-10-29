package main

import (
	"beetle-quest/pkg/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	internalAuthToken = os.Getenv("INTERNAL_AUTH_TOKEN")
)

func main() {
	if internalAuthToken == "" {
		panic("INTERNAL_AUTH_TOKEN is not set")
	}

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
	r.Use(middleware.CheckAuthServiceMiddleware(internalAuthToken))

	basePath := r.Group("/api/v1/user")
	{
		f := func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "EHEHEHEH NOT IMPLEMENTED YET!"})
		}

		basePath.GET("/account/:user_id", f)
		basePath.PATCH("/account/:user_id", f)
		basePath.DELETE("/account/:user_id", f)

		basePath.GET("/:user_id/gacha", f)
		basePath.GET("/:user_id/gacha/:gacha_id", f)
	}

	r.Run()
}
