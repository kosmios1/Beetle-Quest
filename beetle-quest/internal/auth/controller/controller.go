package controller

import (
	"beetle-quest/internal/auth/service"
	"beetle-quest/pkg/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	sessionID := uuid.New().String()

	session := sessions.Default(ctx)
	session.Set("username", loginData.Username)
	session.Set("session_id", sessionID)
	session.Save()

	ctx.JSON(http.StatusOK, gin.H{"message": "logged in as " + session.Get("username").(string)})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if session.Get("session_id") == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	session.Clear()
	session.Save()
	ctx.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
