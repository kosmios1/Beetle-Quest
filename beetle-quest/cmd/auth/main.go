package main

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/middleware"
	"beetle-quest/internal/auth/repository"
	"beetle-quest/internal/auth/service"
	"encoding/hex"
	"os"

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

	internalAuthToken = os.Getenv("INTERNAL_AUTH_TOKEN")
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
		Secure:   true,
		MaxAge:   60 * 20,
		HttpOnly: true, // NOTE: This is set to true to prevent XSS attacks
	})

	return store
}

func main() {
	if internalAuthToken == "" {
		panic("INTERNAL_AUTH_TOKEN is not set")
	}

	// This will connect to redis and return a store object used by the session middleware to store session data
	store := setup_redis_connection()

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("my-session", store))

	{
		cnt := controller.AuthController{
			AuthService: service.AuthService{
				UserRepo: repository.NewUserRepo(),
			},
		}
		r.GET("/logout", cnt.Logout)
		r.POST("/login", cnt.Login)
		r.POST("/register", cnt.Register)
	}

	authorized := r.Group("/api/v1", controller.Proxy)
	authorized.Use(middleware.AuthMiddleware(internalAuthToken))

	r.Run()
}
