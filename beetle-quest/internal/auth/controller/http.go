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

	session := sessions.Default(ctx)
	if session.Get("session_id") != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "already logged in!"})
		ctx.Abort()
		return
	}

	sessionID := utils.GenerateUUID()
	session.Set("user_id", user.UserID.String())
	session.Set("session_id", sessionID.String())
	session.Save()

	ctx.HTML(http.StatusOK, "home.tmpl", gin.H{"UserID": user.UserID.String()})
}

func (c *AuthController) CheckSession(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if session.Get("session_id") == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx.HTML(http.StatusOK, "home.tmpl", gin.H{"UserID": session.Get("user_id")})
}

func (c *AuthController) Authorize(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if session.Get("session_id") == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if session.Get("session_id") == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	session.Clear()
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/")
}
