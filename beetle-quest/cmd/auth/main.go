package main

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/repository"
	internalRepo "beetle-quest/internal/auth/repository"
	"beetle-quest/internal/auth/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// This will connect to redis and return a store object used by the session middleware to store session data
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
	basePath := r.Group("/api/v1/auth")
	{
		cnt := controller.AuthController{
			AuthService: service.AuthService{
				UserRepo:   repository.NewUserRepo(),
				Oauth2Repo: internalRepo.NewOauth2Repo(),
			},
		}
		basePath.GET("/logout", cnt.Logout)
		basePath.POST("/login", cnt.Login)
		basePath.POST("/register", cnt.Register)
		basePath.GET("/oauth2", cnt.Oauth2Callback)

		basePath.GET("/check_session", cnt.CheckSession)
		basePath.Any("/traefik/verify", cnt.Verify)
	}

	r.Run()
}
