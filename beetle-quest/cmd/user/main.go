package main

import (
	"beetle-quest/internal/user/controller"
	"beetle-quest/internal/user/repository"
	"beetle-quest/internal/user/service"
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

	cnt := controller.UserController{
		UserService: service.UserService{
			UserRepo: repository.NewUserRepo(),
		},
	}

	basePath := r.Group("/api/v1/user")
	{
		basePath.GET("/account/:user_id", cnt.GetUserAccountDetails)
		basePath.PATCH("/account/:user_id", cnt.UpdateUserAccountDetails)
		basePath.DELETE("/account/:user_id", cnt.DeleteUserAccount)

		basePath.GET("/:user_id/gacha", cnt.GetUserGachaList)
		basePath.GET("/:user_id/gacha/:gacha_id", cnt.GetUserGachaDetails)
	}

	r.Run()
}
