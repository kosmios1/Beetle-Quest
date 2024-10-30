package controller

import (
	"beetle-quest/internal/user/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service.UserService
}

func (c *UserController) GetUserAccountDetails(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "EHEHEHEH NOT IMPLEMENTED YET!"})
}

func (c *UserController) UpdateUserAccountDetails(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "EHEHEHEH NOT IMPLEMENTED YET!"})
}

func (c *UserController) DeleteUserAccount(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "EHEHEHEH NOT IMPLEMENTED YET!"})
}

func (c *UserController) GetUserGachaList(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "EHEHEHEH NOT IMPLEMENTED YET!"})
}

func (c *UserController) GetUserGachaDetails(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "EHEHEHEH NOT IMPLEMENTED YET!"})
}
