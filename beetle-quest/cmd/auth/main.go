package main

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/middleware"
	"beetle-quest/internal/auth/service"
	repository "beetle-quest/pkg/repositories/impl"
	"net/http"

	"encoding/hex"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

var (
	redisHost       = os.Getenv("REDIS_HOST")
	redisPort       = os.Getenv("REDIS_PORT")
	redisPasswd     = os.Getenv("REDIS_PASSWD")
	redisEncSecret  = os.Getenv("REDIS_ENC_SECRET")
	redisAuthSecret = os.Getenv("REDIS_AUTH_SECRET")
)

func setup_redis_connection() redis.Store {
	if redisHost == "" || redisPort == "" || redisPasswd == "" || redisAuthSecret == "" || redisEncSecret == "" {
		panic("Either REDIS_HOST, REDIS_PORT, REDIS_PASSWD, REDIS_ENC_SECRET or REDIS_AUTH_SECRET is not set")
	}

	auth_secret, err := hex.DecodeString(redisAuthSecret)
	if err != nil {
		panic("Could not decode REDIS_AUTH_SECRET as hex string")
	} else if len(auth_secret) != 64 {
		panic("REDIS_AUTH_SECRET must be 64 bytes long")
	}

	enc_secret, err := hex.DecodeString(redisEncSecret)
	if err != nil {
		panic("Could not decode REDIS_ENC_SECRET as hex string")
	} else if len(enc_secret) != 32 {
		panic("REDIS_ENC_SECRET must be 32 bytes long")
	}

	store, err := redis.NewStore(10, "tcp", redisHost+":"+redisPort, redisPasswd, auth_secret, enc_secret)
	if err != nil {
		panic(err)
	}

	store.Options(sessions.Options{
		Path:     "/",
		Secure:   true,
		HttpOnly: true, // NOTE: This is set to true to prevent XSS attacks
		MaxAge:   int(time.Hour) * 24,
	})

	return store
}

func main() {
	// This will connect to redis and return a store object used by the session middleware to store session data
	store := setup_redis_connection()

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

	r.Use(sessions.Sessions("my-session", store))
	r.LoadHTMLGlob("templates/*")

	// Serve static files
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	r.GET("/static/*filepath", controller.Proxy("http://static-service:8080"))

	basePath := r.Group("/api/v1")
	{
		cnt := controller.AuthController{
			AuthService: service.AuthService{
				UserRepo: repository.NewUserRepo(),
			},
		}
		basePath.GET("/logout", cnt.Logout)
		basePath.POST("/login", cnt.Login)
		basePath.POST("/register", cnt.Register)
		basePath.GET("/check_session", cnt.CheckSession)
	}

	authorized := r.Group("/api/v1")
	authorized.Use(middleware.AuthMiddleware())
	{
		// User routes ===================================================================================================================
		userServiceAddr := "http://user-service:8080"
		userGroup := authorized.Group("/user")
		{
			userGroup.GET("/account/:user_id", controller.Proxy(userServiceAddr))
			userGroup.PATCH("/account/:user_id", controller.Proxy(userServiceAddr))
			userGroup.DELETE("/account/:user_id", controller.Proxy(userServiceAddr))

			userGroup.GET("/:user_id/gacha", controller.Proxy(userServiceAddr))
			userGroup.GET("/:user_id/gacha/:gacha_id", controller.Proxy(userServiceAddr))
		}

		// Gacha routes ==================================================================================================================
		gachaServiceAddr := "http://gacha-service:8080"
		gachaGroup := authorized.Group("/gacha")
		{
			gachaGroup.POST("/roll", controller.Proxy(gachaServiceAddr))
			gachaGroup.GET("/list", controller.Proxy(gachaServiceAddr))
			gachaGroup.GET("/:gacha_id", controller.Proxy(gachaServiceAddr))
		}

		// Market routes =================================================================================================================
		marketServiceAddr := "http://market-service:8080"
		marketGroup := authorized.Group("/market")
		{
			marketGroup.POST("/bugscoin/buy", controller.Proxy(marketServiceAddr))

			auctionGroup := marketGroup.Group("/auction")
			{
				auctionGroup.POST("/create", controller.Proxy(marketServiceAddr))
				auctionGroup.GET("/list", controller.Proxy(marketServiceAddr))
				auctionGroup.GET("/:auction_id", controller.Proxy(marketServiceAddr))
				auctionGroup.DELETE("/:auction_id", controller.Proxy(marketServiceAddr))
				auctionGroup.POST("/:auction_id/bid", controller.Proxy(marketServiceAddr))
			}
		}

		// Admin routes ==================================================================================================================
		adminServiceAddr := "http://admin-service:8080"
		adminGroup := authorized.Group("/admin")
		{
			adminGroup.GET("/monitor", controller.Proxy(adminServiceAddr))
			adminGroup.GET("/user/list", controller.Proxy(adminServiceAddr))
			adminGroup.POST("/gacha", controller.Proxy(adminServiceAddr))
			adminGroup.PATCH("/gacha/:gacha_id", controller.Proxy(adminServiceAddr))
			adminGroup.POST("/report/issue", controller.Proxy(adminServiceAddr))
		}
	}
	r.Run()
}
