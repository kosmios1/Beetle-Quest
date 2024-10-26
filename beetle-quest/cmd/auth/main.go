package main

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/middleware"
	"beetle-quest/internal/auth/service"
	"encoding/hex"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func setup_redis_connection() *redis.Store {
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_passwd := os.Getenv("REDIS_PASSWD")
	redis_enc_secret := os.Getenv("REDIS_ENC_SECRET")
	redis_auth_secret := os.Getenv("REDIS_AUTH_SECRET")

	if redis_host == "" || redis_port == "" || redis_passwd == "" || redis_auth_secret == "" || redis_enc_secret == "" {
		panic("Either REDIS_HOST, REDIS_PORT, REDIS_PASSWD, REDIS_ENC_SECRET or REDIS_AUTH_SECRET is not set")
	}

	auth_secret, err := hex.DecodeString(redis_auth_secret)
	if err != nil {
		panic("Could not decode REDIS_AUTH_SECRET as hex string")
	} else if len(auth_secret) != 64 {
		panic("REDIS_AUTH_SECRET must be 64 bytes long")
	}

	enc_secret, err := hex.DecodeString(redis_enc_secret)
	if err != nil {
		panic("Could not decode REDIS_ENC_SECRET as hex string")
	} else if len(enc_secret) != 32 {
		panic("REDIS_ENC_SECRET must be 32 bytes long")
	}

	store, err := redis.NewStore(10, "tcp", redis_host+":"+redis_port, redis_passwd, auth_secret, enc_secret)
	if err != nil {
		panic(err)
	}

	store.Options(sessions.Options{
		Secure: true,
	})

	return &store
}

func main() {
	internalAuthToken, ok := os.LookupEnv("INTERNAL_AUTH_TOKEN")
	if !ok {
		panic("INTERNAL_AUTH_TOKEN is not set")
	}

	store := setup_redis_connection()

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("user_sessions", *store))

	{
		cnt := controller.AuthController{
			AuthService: service.AuthService{
				UserRepo: nil, // NOTE: setup user repository
			},
		}
		r.GET("/logout", cnt.Logout)
		r.POST("/login", cnt.Login)
		r.POST("/register", cnt.Register)
	}

	authorized := r.Group("/api/v1", controller.Proxy)
	authorized.Use(middleware.AuthMiddleware(internalAuthToken))
	// {
	// 	// NOTE: Can this be generalized? at least the controller.Proxy?
	// 	authorized.Any("/user/*", controller.Proxy)
	// 	authorized.Any("/gacha/*", controller.Proxy)
	// 	authorized.Any("/market/*", controller.Proxy)
	// 	authorized.Any("/auction/*", controller.Proxy)
	// 	authorized.Any("/admin/*", controller.Proxy)
	// 	authorized.Any("/report/*", controller.Proxy)
	// }

	r.Run()
}
