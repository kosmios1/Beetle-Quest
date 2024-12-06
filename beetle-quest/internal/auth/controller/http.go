package controller

import (
	"beetle-quest/internal/auth/service"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
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
			oauth2.Refreshing,
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

	o2srv.SetResponseTokenHandler(func(w http.ResponseWriter, data map[string]interface{}, header http.Header, statusCode ...int) error {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")

		for key := range header {
			w.Header().Set(key, header.Get(key))
		}

		status := http.StatusOK
		if len(statusCode) > 0 && statusCode[0] > 0 {
			status = statusCode[0]
		}

		if _, ok := data["access_token"]; ok {
			accessTokenClaims, err := utils.VerifyJWTToken(data["access_token"].(string), jwtSecretKey)
			if err != nil {
				return err
			}
			claims := jwt.MapClaims{
				"sub": accessTokenClaims["sub"],
				"iss": "Beetle Quest",
				"aud": "beetle-quest",
				"iat": time.Now().Unix(),
				"nbf": time.Now().Unix(),
				"exp": time.Now().Add(time.Hour).Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			idTok, err := token.SignedString(jwtSecretKey)
			if err != nil {
				return err
			}

			data["id_token"] = idTok
		}

		w.WriteHeader(status)
		return json.NewEncoder(w).Encode(data)
	})

	o2srv.SetUserAuthorizationHandler(cnt.userAuthorizationHandler)
	o2srv.SetAuthorizeScopeHandler(cnt.authorizeScopeHandler)

	return cnt
}

func (c *AuthController) AuthenticationPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "loginPage.tmpl", gin.H{})
}

func (c *AuthController) AuthorizePage(ctx *gin.Context) {
	store, err := session.Start(ctx.Request.Context(), ctx.Writer, ctx.Request)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		ctx.Redirect(http.StatusFound, "/api/v1/auth/authPage")
		return
	}
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
		case models.ErrUsernameOrEmailAlreadyExists, models.ErrInvalidUsernameOrPassOrEmail, models.ErrInvalidUUID:
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

	user, err := c.srv.Login(loginData.Username, loginData.Password)
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInvalidPassword, models.ErrInvalidUUID:
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

	store, err := session.Start(ctx.Request.Context(), ctx.Writer, ctx.Request)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	store.Set("LoggedInUserID", user.UserID.String())
	store.Set("IsAdmin", false)
	if err := store.Save(); err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.Redirect(http.StatusFound, "/api/v1/auth/authorizePage")
}

func (c *AuthController) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) == 2 {
		if err := c.o2mng.RemoveAccessToken(ctx, bearerToken[1]); err != nil {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": models.ErrInternalServerError})
			ctx.Abort()
			return
		}
	}

	if err := session.Destroy(ctx.Request.Context(), ctx.Writer, ctx.Request); err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": models.ErrInternalServerError})
		ctx.Abort()
		return
	}

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
	store, err := session.Start(ctx.Request.Context(), ctx.Writer, ctx.Request)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
		ctx.Request.Form = form
	}

	store.Delete("ReturnUri")
	_ = store.Save()

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
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			_ = r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		_ = store.Save()

		w.Header().Set("Location", "/api/v1/auth/authPage")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	_ = store.Save()
	return
}

func (c *AuthController) authorizeScopeHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	scopes := r.FormValue("scope")
	for _, scope := range strings.Split(scopes, ", ") {
		if scope == "admin" {
			store, err := session.Start(r.Context(), w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return "", err
			}

			isAdmin, ok := store.Get("IsAdmin")
			if !ok || !isAdmin.(bool) {
				return "", o2errors.ErrUnauthorizedClient
			}
			break
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

	admin, err := c.srv.AdminLogin(loginData.AdminID, loginData.Password, loginData.OtpCode)
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInvalidAdminIDOrPassOrOTOP, models.ErrInvalidUUID:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	store, err := session.Start(ctx.Request.Context(), ctx.Writer, ctx.Request)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	store.Set("LoggedInUserID", admin.AdminId.String())
	store.Set("IsAdmin", true)
	if err := store.Save(); err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.Redirect(http.StatusFound, "/api/v1/auth/authorizePage")
}
