package controller

import (
	"beetle-quest/internal/auth/service"
	"beetle-quest/pkg/models"
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
		return
	}

	if err := c.AuthService.Register(registerData.Email, registerData.Username, registerData.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var loginData models.LoginRequest
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email or password is missing"})
		return
	}

	if err := c.AuthService.Login(loginData.Username, loginData.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	session := sessions.Default(ctx)
	session.Set("username", loginData.Username)
	session.Save()

	ctx.JSON(http.StatusOK, gin.H{"message": "logged in as " + loginData.Username})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("username")
	session.Save()
	ctx.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
