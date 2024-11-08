package main

import (
	"beetle-quest/internal/user/controller"
	"beetle-quest/internal/user/service"
	repository "beetle-quest/pkg/repositories/impl"

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

	cnt := controller.UserController{
		UserService: service.UserService{
			UserRepo: repository.NewUserRepo(),
		},
	}

	basePath := r.Group("/api/v1/user")
	{
		accountGroup := basePath.Group("/account")
		{
			accountGroup.GET("/:user_id", cnt.GetUserAccountDetails)
			accountGroup.PATCH("/:user_id", cnt.UpdateUserAccountDetails)
			accountGroup.DELETE("/:user_id", cnt.DeleteUserAccount)
		}

		basePath.GET("/:user_id/gacha", cnt.GetUserGachaList)
		basePath.GET("/:user_id/gacha/:gacha_id", cnt.GetUserGachaDetails)
	}

	r.Run()
}
