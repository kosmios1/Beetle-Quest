package controller

import (
	"beetle-quest/internal/auth/service"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service.AuthService
}

func (c *AuthController) Register(ctx *gin.Context) {
	var registerData models.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email or password is missing"})
		ctx.Abort()
		return
	}

	if err := c.AuthService.Register(registerData.Email, registerData.Username, registerData.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var loginData models.LoginRequest
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email or password is missing"})
		ctx.Abort()
		return
	}

	user, err := c.AuthService.Login(loginData.Username, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		ctx.Abort()
		return
	}

	sessionID := utils.GenerateUUID()
	session := sessions.Default(ctx)
	session.Set("user_id", user.UserID.String())
	session.Set("session_id", sessionID.String())
	session.Save()

	// NOTE: This works because of http.StatusFound and not http.StatusPermanentRedirect
	// location := url.URL{Path: "/api/v1/user/account/" + user.UserID.String()}
	// ctx.Redirect(http.StatusFound, location.RequestURI())
	ctx.JSON(http.StatusOK, gin.H{"message": "logged in"})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if session.Get("session_id") == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		ctx.Abort()
		return
	}

	session.Clear()
	session.Save()
	ctx.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
