package controller

import (
	"beetle-quest/internal/auth/service"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"log"
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
	ctx.SetCookie("access_token", tokenString, int(maxAge), "/", "", false, true)
	ctx.Redirect(http.StatusFound, "/api/v1/auth/check_session")
}

func (c *AuthController) Logout(ctx *gin.Context) {
	token, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.Redirect(http.StatusFound, "/static/")
		return
	}

	if ok := c.srv.RevokeToken(token); !ok {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": models.ErrInternalServerError})
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
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": models.ErrInternalServerError})
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
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": models.ErrInternalServerError})
		ctx.Abort()
		return
	}

	if claims["scope"] == "user" {
		ctx.HTML(http.StatusOK, "home.tmpl", gin.H{"UserID": claims["sub"]})
	} else {
		ctx.HTML(http.StatusOK, "adminHome.tmpl", gin.H{"AdminID": claims["sub"]})
	}
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
	ctx.SetCookie("access_token", tokenString, int(maxAge), "/", "", false, true)
	ctx.Redirect(http.StatusFound, "/api/v1/auth/check_session")
}
