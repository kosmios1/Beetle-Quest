package controller

import (
	"beetle-quest/internal/user/service"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service.UserService
}

func (c *UserController) GetUserAccountDetails(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	parsedUserID, err := utils.ParseUUID(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	user, err := c.UserService.GetUserAccountDetails(parsedUserID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	// TODO: Get user's gacha list
	// TODO: Get user's transaction history

	ctx.HTML(http.StatusOK, "userInfo.tmpl", models.GetUserAccountDetailsTemplatesData{
		UserID:       user.UserID.String(),
		Username:     user.Username,
		Email:        user.Email,
		Currency:     user.Currency,
		Gachas:       []models.Gacha{},
		Transactions: []models.Transaction{},
	})
}

func (c *UserController) UpdateUserAccountDetails(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	parsedUserID, err := utils.ParseUUID(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	var req models.UpdateUserAccountDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	err = c.UserService.UpdateUserAccountDetails(parsedUserID, req.Email, req.Username, req.OldPassword, req.NewPassword)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{
		"Message": "User account updated successfully",
	})
}

func (c *UserController) DeleteUserAccount(ctx *gin.Context) {

	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	parsedUserID, err := utils.ParseUUID(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	err = c.UserService.DeleteUserAccount(parsedUserID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User account not found"})
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{
		"Message": "User account deleted successfully!",
	})
}

func (c *UserController) GetUserGachaList(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "EHEHEHEH NOT IMPLEMENTED YET!"})
}

func (c *UserController) GetUserGachaDetails(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "EHEHEHEH NOT IMPLEMENTED YET!"})
}
