package controller

import (
	"beetle-quest/internal/auth/service"
	"beetle-quest/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	srv *service.AuthService
}

func NewAuthController(srv *service.AuthService) *AuthController {
	return &AuthController{
		srv: srv,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var registerData models.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerData); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "email or password is missing"})
		ctx.Abort()
		return
	}

	if err := c.srv.Register(registerData.Email, registerData.Username, registerData.Password); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err})
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "User registered successfully!"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var loginData models.LoginRequest
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "username or password is missing"})
		ctx.Abort()
		return
	}

	user, err := c.srv.Login(loginData.Username, loginData.Password)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err})
		ctx.Abort()
		return
	}

	url, err := c.srv.MakeAuthRequest(string(user.UserID.String()))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": "internal server error"})
		ctx.Abort()
		return
	}

	ctx.Redirect(http.StatusFound, url)
}

func (c *AuthController) Oauth2Callback(ctx *gin.Context) {
	if err := ctx.Request.ParseForm(); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "invalid request"})
		ctx.Abort()
		return
	}

	state := ctx.Request.Form.Get("state")
	// TODO: How to validate state?
	if state == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "state invalid!"})
		ctx.Abort()
		return
	}

	code := ctx.Request.Form.Get("code")
	if code == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "code not found!"})
		ctx.Abort()
		return
	}

	token, claims, err := c.srv.ExchangeCodeForToken(code)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": "internal server error"})
		ctx.Abort()
		return
	}

	maxAge := int(token.Expiry.Sub(time.Now()).Seconds())
	ctx.SetCookie("access_token", token.AccessToken, maxAge, "/", "", false, true)

	ctx.HTML(http.StatusOK, "home.tmpl", gin.H{"UserID": claims["sub"]})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	token, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.Redirect(http.StatusFound, "/static/")
		return
	}

	if ok := c.srv.RevokeToken(token); !ok {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": "internal server error"})
		ctx.Abort()
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "", true, true)
	ctx.Redirect(http.StatusFound, "/static/")
}

func (c *AuthController) Verify(ctx *gin.Context) {
	token, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.Redirect(http.StatusFound, "/static/")
		return
	}

	if _, ok := c.srv.VerifyToken(token); !ok {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": "internal server error"})
		ctx.Abort()
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *AuthController) CheckSession(ctx *gin.Context) {
	token, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	claims, ok := c.srv.VerifyToken(token)
	if !ok {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": "internal server error"})
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "home.tmpl", gin.H{"UserID": claims["sub"]})
}
