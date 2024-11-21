package controller

import (
	"beetle-quest/internal/user/service"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	srv *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service,
	}
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

	user, err := c.srv.GetUserAccountDetails(parsedUserID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err})
		ctx.Abort()
		return
	}

	gachas, err := c.srv.GetUserGachaList(user.UserID.String())
	if err != nil {
		if err == models.ErrInternalServerError {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
			return
		} else {
			gachas = []models.Gacha{}
		}
	}

	transactions := c.srv.GetUserTransactionHistory(userID)

	var transactionViews []models.TransactionView
	for _, transaction := range transactions {
		transactionViews = append(transactionViews, models.TransactionView{
			TransactionID:   transaction.TransactionID.String(),
			TransactionType: transaction.TransactionType.String(),
			UserID:          transaction.UserID.String(),
			Amount:          transaction.Amount,
			DateTime:        transaction.DateTime,
			EventType:       transaction.EventType.String(),
			EventID:         transaction.EventID.String(),
		})
	}

	ctx.HTML(http.StatusOK, "userInfo.tmpl", models.GetUserAccountDetailsTemplatesData{
		UserID:          user.UserID.String(),
		Username:        user.Username,
		Email:           user.Email,
		Currency:        user.Currency,
		GachaList:       gachas,
		TransactionList: transactionViews,
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

	err = c.srv.UpdateUserAccountDetails(parsedUserID, req.Email, req.Username, req.OldPassword, req.NewPassword)
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

	err = c.srv.DeleteUserAccount(parsedUserID, password)
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
			return
		case models.ErrUserNotFound, models.ErrRetalationGachaUserNotFound, models.ErrUserTransactionNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
			return
		}
		panic("unreachable code")
	}

	ctx.Redirect(http.StatusSeeOther, "/api/v1/auth/logout")
}

// Internal API ==========================================================================

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users := c.srv.GetAllUsers()

	var data models.GetAllUsersDataResponse = models.GetAllUsersDataResponse{
		UserList: users,
	}
	ctx.JSON(http.StatusOK, data)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req models.CreateUserData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Wrong inputs passed to the request!"})
		return
	}

	if ok := c.srv.Create(req.Email, req.Username, req.HashedPassword, req.Currency); !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "internal server error!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "User created successfully"})
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	var req models.User
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Wrong inputs passed to the request!"})
		return
	}

	if ok := c.srv.Update(&req); !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "internal server error!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "User updated successfully"})
}

func (c *UserController) FindByID(ctx *gin.Context) {
	var req models.FindUserByIDData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Wrong inputs passed to the request!"})
		ctx.Abort()
		return
	}

	user, exits := c.srv.FindByID(req.UserID)
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

	user, exits := c.srv.FindByUsername(req.Username)
	if !exits {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "User not found!"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
