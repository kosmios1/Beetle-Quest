package controller

import (
	"beetle-quest/internal/auth/service"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"encoding/hex"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/go-oauth2/oauth2/v4"
	o2errors "github.com/go-oauth2/oauth2/v4/errors"
	o2manage "github.com/go-oauth2/oauth2/v4/manage"
	o2models "github.com/go-oauth2/oauth2/v4/models"
	o2server "github.com/go-oauth2/oauth2/v4/server"
	o2store "github.com/go-oauth2/oauth2/v4/store"

	accessToken "beetle-quest/internal/auth/service/jwtAccessToken"
	store "beetle-quest/internal/auth/service/oauth2/storage"
)

var (
	jwtSecretKey = utils.PanicIfError[[]byte](hex.DecodeString(utils.FindEnv("JWT_SECRET_KEY")))
)

type AuthController struct {
	srv *service.AuthService

	o2mng *o2manage.Manager
	o2srv *o2server.Server
}

func NewAuthController(srv *service.AuthService) *AuthController {
	manager := o2manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(o2manage.DefaultAuthorizeCodeTokenCfg)

	manager.MapTokenStorage(store.GetTokenStorage())
	manager.MapAccessGenerate(accessToken.NewJWTAccessGenerate("", jwtSecretKey, jwt.SigningMethodHS512))

	clientStore := o2store.NewClientStore()
	if err := clientStore.Set("beetle-quest", &o2models.Client{ID: "beetle-quest"}); err != nil {
		log.Fatal("[FATAL] Could not setup oauth2 client storage!")
	}
	manager.MapClientStorage(clientStore)

	o2srv := o2server.NewServer(&o2server.Config{
		TokenType:            "Bearer",
		AllowedResponseTypes: []oauth2.ResponseType{oauth2.Code},
		AllowedGrantTypes: []oauth2.GrantType{
			oauth2.AuthorizationCode,
		},
		AllowedCodeChallengeMethods: []oauth2.CodeChallengeMethod{
			oauth2.CodeChallengeS256,
		},
		AllowGetAccessRequest: false,
		ForcePKCE:             true,
	}, manager)
	o2srv.SetClientInfoHandler(o2server.ClientFormHandler)

	o2srv.SetInternalErrorHandler(func(err error) (re *o2errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	o2srv.SetResponseErrorHandler(func(re *o2errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	cnt := &AuthController{
		srv: srv,

		o2mng: manager,
		o2srv: o2srv,
	}

	o2srv.SetUserAuthorizationHandler(cnt.userAuthorizationHandler)
	o2srv.SetAuthorizeScopeHandler(cnt.authorizeScopeHandler)

	return cnt
}

func (c *AuthController) AuthenticationPage(ctx *gin.Context) {
	redirect, _ := ctx.GetQuery("redirect")
	if redirect != "" {
		decodedRedirect, err := url.QueryUnescape(redirect)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode redirect parameter"})
			return
		}
		parsedURL, err := url.Parse(decodedRedirect)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse redirect URL"})
			return
		}
		redirectURI := parsedURL.Query().Get("redirect_uri")
		if redirectURI == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing redirect_uri in the decoded URL"})
			return
		}

		secPolicy := ctx.GetHeader("Content-Security-Policy")
		ctx.Header("Content-Security-Policy", secPolicy+"; connect-src 'self' "+redirectURI+";")
	}
	ctx.HTML(http.StatusOK, "loginPage.tmpl", gin.H{"Redirect": redirect})
}

func (c *AuthController) AuthorizePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "authorizePage.tmpl", gin.H{})
}

func (c *AuthController) TokenPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "tokenPage.tmpl", gin.H{
		"State": ctx.Query("state"),
		"Code":  ctx.Query("code"),
	})
}

func (c *AuthController) Register(ctx *gin.Context) {
	var registerData models.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerData); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "email or password is missing"})
		ctx.Abort()
		return
	}

	if err := c.srv.Register(registerData.Email, registerData.Username, registerData.Password); err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrUsernameOrEmailAlreadyExists, models.ErrInvalidUsernameOrPassOrEmail:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return

		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	ctx.HTML(http.StatusCreated, "successMsg.tmpl", gin.H{"Message": "User registered successfully!"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var loginData models.LoginRequest
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	token, tokenString, err := c.srv.Login(loginData.Username, loginData.Password)
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInvalidPassword:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrUserNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	maxAge := token.Claims.(*utils.CustomClaims).ExpiresAt - time.Now().Unix()
	ctx.SetCookie("identity_token", tokenString, int(maxAge), "/", "", false, true)

	if loginData.Redirect != "" {
		ctx.Redirect(http.StatusFound, loginData.Redirect)
	} else {
		ctx.Redirect(http.StatusFound, "/api/v1/auth/authorizePage")
	}
}

func (c *AuthController) Logout(ctx *gin.Context) {
	token, err := ctx.Cookie("identity_token")
	if err != nil {
		ctx.Redirect(http.StatusFound, "/static/")
		return
	}

	if ok := c.srv.RevokeToken(token); !ok {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": models.ErrInternalServerError})
		ctx.Abort()
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": models.ErrInternalServerError})
		ctx.Abort()
		return
	}

	if err := c.o2mng.RemoveAccessToken(ctx, bearerToken[1]); err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": models.ErrInternalServerError})
		ctx.Abort()
		return
	}

	ctx.SetCookie("identity_token", "", -1, "/", "", true, true)
	ctx.Redirect(http.StatusFound, "/static/")
}

func (c *AuthController) Verify(ctx *gin.Context) {
	_, err := c.o2srv.ValidationBearerToken(ctx.Request)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/api/v1/auth/authPage")
		ctx.Abort()
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *AuthController) CheckSession(ctx *gin.Context) {
	tokenInfo, err := c.o2srv.ValidationBearerToken(ctx.Request)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/api/v1/auth/authPage")
		ctx.Abort()
		return
	}

	isAdmin := false
	scopes := tokenInfo.GetScope()
	for _, scope := range strings.Split(scopes, ", ") {
		if scope == "admin" {
			isAdmin = true
			break
		}
	}

	if isAdmin == false {
		ctx.HTML(http.StatusOK, "home.tmpl", gin.H{"UserID": tokenInfo.GetUserID()})
	} else {
		ctx.HTML(http.StatusOK, "adminHome.tmpl", gin.H{"AdminID": tokenInfo.GetUserID()})
	}
}

// Oauth ==============================================================================================================-

func (c *AuthController) OauthAuthorize(ctx *gin.Context) {
	if err := c.o2srv.HandleAuthorizeRequest(ctx.Writer, ctx.Request); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}
}

func (c *AuthController) OauthToken(ctx *gin.Context) {
	if err := c.o2srv.HandleTokenRequest(ctx.Writer, ctx.Request); err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}
}

func (c *AuthController) userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	cookie, err := r.Cookie("identity_token")
	if err != nil {
		redirect := "/api/v1/auth/authPage?redirect=" + url.QueryEscape(r.URL.RequestURI())
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}

	claims, ok := c.srv.VerifyToken(cookie.Value)
	if !ok {
		redirect := "/api/v1/auth/authPage?redirect=" + url.QueryEscape(r.URL.RequestURI())
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}

	return claims["sub"].(string), nil
}

func (c *AuthController) authorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	scopes := r.FormValue("scope")
	for _, scope := range strings.Split(scopes, ", ") {
		if scope == "admin" {
			cookie, err := r.Cookie("identity_token")
			if err != nil {
				redirect := "/api/v1/auth/authPage?redirect=" + url.QueryEscape(r.URL.RequestURI())
				http.Redirect(w, r, redirect, http.StatusFound)
				return "", err
			}

			claims, ok := c.srv.VerifyToken(cookie.Value)
			if !ok {
				redirect := "/api/v1/auth/authPage?redirect=" + url.QueryEscape(r.URL.RequestURI())
				http.Redirect(w, r, redirect, http.StatusFound)
				return "", o2errors.ErrServerError
			}

			if claims["is_admin"].(bool) == false {
				return "", o2errors.ErrInvalidScope
			}
		}
	}
	return scopes, nil
}

// Admin ==============================================================================================================

func (c *AuthController) AdminLogin(ctx *gin.Context) {
	var loginData models.AdminLoginRequest
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	token, tokenString, err := c.srv.AdminLogin(loginData.AdminID, loginData.Password, loginData.OtpCode)
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInvalidAdminIDOrPassOrOTOP:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	maxAge := token.Claims.(*utils.CustomClaims).ExpiresAt - time.Now().Unix()
	ctx.SetCookie("identity_token", tokenString, int(maxAge), "/", "", false, true)

	if loginData.Redirect != "" {
		ctx.Redirect(http.StatusFound, loginData.Redirect)
	} else {
		ctx.Redirect(http.StatusFound, "/api/v1/auth/authorizePage")
	}
}
