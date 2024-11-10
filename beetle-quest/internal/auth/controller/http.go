package controller

import (
	"beetle-quest/internal/auth/service"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"context"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service.AuthService
}

var (
	jwtSecretKey = utils.PanicIfError[[]byte](hex.DecodeString(os.Getenv("JWT_SECRET_KEY")))
)

func (c *AuthController) Register(ctx *gin.Context) {
	var registerData models.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerData); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "email or password is missing"})
		ctx.Abort()
		return
	}

	if err := c.AuthService.Register(registerData.Email, registerData.Username, registerData.Password); err != nil {
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

	user, err := c.AuthService.Login(loginData.Username, loginData.Password)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err})
		ctx.Abort()
		return
	}

	state, err := utils.GenerateRandomSalt(32)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": "internal server error"})
		ctx.Abort()
		return
	}
	stateHex := hex.EncodeToString(state)

	url := c.AuthCodeURL(stateHex, user.UserID.String())
	if url == "" {
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
	}

	token, err := c.Exchange(context.Background(), code)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": "internal server error"})
		ctx.Abort()
		return
	}

	claims, err := utils.VerifyJWTToken(token.AccessToken, jwtSecretKey)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": "internal server error"})
		ctx.Abort()
		return
	}

	maxAge := int(token.Expiry.Sub(time.Now()).Seconds())
	// TODO: Secure when https
	ctx.SetCookie("access_token", token.AccessToken, maxAge, "/", "", false, true)
	ctx.HTML(http.StatusOK, "home.tmpl", gin.H{"UserID": claims["sub"]})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	token, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.Redirect(http.StatusFound, "/static/")
		return
	}

	if _, err := c.RevokeToken(token); err != nil {
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

	resp, err := c.VerifyToken(token)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": "internal server error"})
		ctx.Abort()
		return
	}

	if resp.StatusCode != http.StatusOK {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Status(resp.StatusCode)
}

func (c *AuthController) CheckSession(ctx *gin.Context) {
	token, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	resp, err := c.VerifyToken(token)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": "internal server error"})
		ctx.Abort()
		return
	}

	if resp.StatusCode != http.StatusOK {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, err := utils.VerifyJWTToken(token, jwtSecretKey)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "home.tmpl", gin.H{"UserID": claims["sub"]})
}
