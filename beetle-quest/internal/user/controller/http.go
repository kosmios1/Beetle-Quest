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
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "No User ID has been provided!"})
		ctx.Abort()
		return
	}

	parsedUserID, err := utils.ParseUUID(userID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "Invalid User ID for this session!"})
		ctx.Abort()
		return
	}

	user, err := c.UserService.GetUserAccountDetails(parsedUserID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err})
		ctx.Abort()
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
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "No User ID has been provided!"})
		ctx.Abort()
		return
	}

	parsedUserID, err := utils.ParseUUID(userID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "Invalid User ID for this session!"})
		ctx.Abort()
		return
	}

	var req models.UpdateUserAccountDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "Wrong inputs passed to the request!"})
		ctx.Abort()
		return
	}

	err = c.UserService.UpdateUserAccountDetails(parsedUserID, req.Email, req.Username, req.OldPassword, req.NewPassword)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err})
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{
		"Message": "User account updated successfully",
	})
}

func (c *UserController) DeleteUserAccount(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "No User ID has been provided!"})
		ctx.Abort()
		return
	}

	parsedUserID, err := utils.ParseUUID(userID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "Invalid User ID for this session!"})
		ctx.Abort()
		return
	}

	password, ok := ctx.GetQuery("password")
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "No password inserted!"})
		ctx.Abort()
		return
	}

	err = c.UserService.DeleteUserAccount(parsedUserID, password)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err})
		ctx.Abort()
		return
	}

	// TODO: Invalidate user's token

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

// Internal API ==========================================================================

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req models.CreateUserData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Wrong inputs passed to the request!"})
		return
	}

	if ok := c.UserService.Create(req.Email, req.Username, req.HashedPassword, req.Currency); !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "internal server error!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "User created successfully"})
}

func (c *UserController) FindByID(ctx *gin.Context) {
	var req models.FindUserByIDData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Wrong inputs passed to the request!"})
		ctx.Abort()
		return
	}

	id, err := utils.ParseUUID(req.UserID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Invalid User ID!"})
		return
	}

	user, exits := c.UserService.FindByID(id)
	if !exits {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "User not found!"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) FindByUsername(ctx *gin.Context) {
	var req models.FindUserByUsernameData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Wrong inputs passed to the request!"})
		ctx.Abort()
		return
	}

	user, exits := c.UserService.FindByUsername(req.Username)
	if !exits {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "User not found!"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
