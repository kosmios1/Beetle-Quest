package main

import (
	"beetle-quest/pkg/utils"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

var (
	clientID     = os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("OAUTH2_CLIENT_SECRET")
	clientDomain = os.Getenv("OAUTH2_CLIENT_DOMAIN")

	jwtKeySecret = utils.PanicIfError[[]byte](hex.DecodeString(os.Getenv("JWT_KEY_SECRET")))
)

func main() {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// Set up JWT access token generation
	// kid : key id, used to distinguish multiple keys
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", jwtKeySecret, jwt.SigningMethodHS512))

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

	r := gin.Default()
	r.Any("/oauth2/token", func(c *gin.Context) {
		if err := srv.HandleTokenRequest(c.Writer, c.Request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	})

	r.Any("/oauth2/authorize", func(c *gin.Context) {
		if err := srv.HandleAuthorizeRequest(c.Writer, c.Request); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
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
