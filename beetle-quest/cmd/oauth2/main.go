package main

import (
	"beetle-quest/pkg/utils"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"

	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"

	myjwtAccess "beetle-quest/cmd/oauth2/jwtAccess"
)

var (
	clientID     string = utils.FindEnv("OAUTH2_CLIENT_ID")
	clientSecret string = utils.FindEnv("OAUTH2_CLIENT_SECRET")
	clientDomain string = utils.FindEnv("OAUTH2_CLIENT_DOMAIN")

	jwtKeySecret []byte = utils.PanicIfError[[]byte](hex.DecodeString(utils.FindEnv("JWT_SECRET_KEY")))

	redisHost     string = utils.FindEnv("REDIS_HOST")
	redisPort     string = utils.FindEnv("REDIS_PORT")
	redisPassword string = utils.FindEnv("REDIS_PASSWORD")
	redisUsername string = utils.FindEnv("REDIS_USERNAME")
	redisDB       int    = utils.PanicIfError[int](strconv.Atoi(utils.FindEnv("REDIS_DB")))
)

func main() {
	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		DB:              redisDB,
		Addr:            redisHost + ":" + redisPort,
		Username:        redisUsername,
		Password:        redisPassword,
		MinRetryBackoff: time.Second * 5,
		MaxRetryBackoff: time.Minute * 2,
	}))
	manager.MapAccessGenerate(myjwtAccess.NewJWTAccessGenerate("", jwtKeySecret, jwt.SigningMethodHS512))

	clientDomain = "" // TODO: remove this line when the client domain is set up
	clientStore := store.NewClientStore()
	clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: clientDomain,
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetClientScopeHandler(authorizeScopeHandler)

	r := gin.Default()
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

	r.Any("/oauth2/authorize", func(ctx *gin.Context) {
		if err := srv.HandleAuthorizeRequest(ctx.Writer, ctx.Request); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	})

	r.Any("/oauth2/token", func(ctx *gin.Context) {
		if err := srv.HandleTokenRequest(ctx.Writer, ctx.Request); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	})

	r.POST("/oauth2/token/revoke", func(ctx *gin.Context) {
		token := ctx.PostForm("token")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "invalid_request"})
			return
		}

		err := srv.Manager.RemoveAccessToken(ctx, token)
		if err != nil {
			if err == errors.ErrInvalidAccessToken {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "invalid_token"})
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "server_error"})
			}
			return
		}

		err = srv.Manager.RemoveRefreshToken(ctx, token)
		if err != nil && err != errors.ErrInvalidRefreshToken {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "server_error"})
			return
		}

		ctx.Status(http.StatusOK)
	})

	r.POST("/oauth2/token/verify", func(ctx *gin.Context) {
		token := ctx.PostForm("token")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "invalid_request"})
			return
		}

		_, err := srv.Manager.LoadAccessToken(ctx, token)
		if err != nil {
			if err == errors.ErrInvalidAccessToken {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "invalid_token"})
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "server_error"})
			}
			return
		}

		ctx.Status(http.StatusOK)
	})

	r.Run()
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	userID = r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return "", err
	}

	return userID, nil
}

func authorizeScopeHandler(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error) {
	if tgr.Scope == "user" || tgr.Scope == "admin" {
		return true, nil
	}
	return false, errors.New("Invalid scope")
}
