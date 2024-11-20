package controller

import (
	"beetle-quest/internal/admin/service"
	"beetle-quest/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	srv *service.AdminService
}

func NewAdminController(srv *service.AdminService) *AdminController {
	return &AdminController{
		srv: srv,
	}
}

func (c *AdminController) FindByID(ctx *gin.Context) {
	var data models.FindAdminByIDData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	admin, exists := c.srv.FindByID(data.AdminID)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}
	ctx.JSON(http.StatusOK, admin)
}

// User controllers =================================================

func (c *AdminController) GetAllUsers(ctx *gin.Context) {
	users, err := c.srv.GetAllUsers()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"error": "Internal server error!"})
		ctx.Abort()
		return
	}

	// ctx.HTML(http.StatusOK, "userList.tmpl", gin.H{"UserList": users})
	ctx.JSON(http.StatusOK, gin.H{"UserList": users})
}

func (c *AdminController) GetUserProfile(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	if userId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	user, err := c.srv.FindUserByID(userId)
	if err != nil {
		if err == models.ErrInternalServerError {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
		} else {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
		}
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *AdminController) UpdateUserProfile(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	if userId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "user_id is required"})
		ctx.Abort()
		return
	}

	var data models.AdminUpdateUserAccount
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "user_id is required"})
		ctx.Abort()
		return
	}

	if ok := c.srv.UpdateUserProfile(userId, &data); !ok {
		ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"error": "Internal server error!"})
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "User profile updated successfully!"})
}

func (c *AdminController) GetUserTransactionHistory(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	if userId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "user_id is required"})
		ctx.Abort()
		return
	}

	transactions, exists := c.srv.GetUserTransactionHistory(userId)
	if !exists {
		ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": "User not found"})
		ctx.Abort()
		return
	}

	// ctx.HTML(http.StatusOK, "userTransactionHistory.tmpl", gin.H{"TransactionList": transactions})
	ctx.JSON(http.StatusOK, gin.H{"TransactionList": transactions})
}

func (c *AdminController) GetUserAuctionList(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	if userId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "user_id is required"})
		ctx.Abort()
		return
	}

	auctionList, exists := c.srv.GetUserAuctionList(userId)
	if !exists {
		ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": "User not found"})
		ctx.Abort()
		return
	}

	// ctx.HTML(http.StatusOK, "userMarketHistory.tmpl", gin.H{"MarketHistory": marketHistory})
	ctx.JSON(http.StatusOK, gin.H{"AuctionList": auctionList})
}

// Gacha controllers =================================================

func (c *AdminController) AddGacha(ctx *gin.Context) {
	var data models.AdminAddGachaRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	if err := c.srv.AddGacha(&data); err != nil {
		if err == models.ErrGachaAlreadyExists {
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
		} else {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
		}
		return
	}
	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Gacha added successfully!"})
}

func (cnt *AdminController) DeleteGacha(ctx *gin.Context) {
	gachaId := ctx.Param("gacha_id")
	if gachaId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "gacha_id is required"})
		ctx.Abort()
		return
	}

	if err := cnt.srv.DeleteGacha(gachaId); err != nil {
		if err == models.ErrGachaNotFound {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
		} else {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
		}
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Gacha deleted successfully!"})
}

func (cnt *AdminController) UpdateGacha(ctx *gin.Context) {
	gachaId := ctx.Param("gacha_id")
	if gachaId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "gacha_id is required"})
		ctx.Abort()
		return
	}

	var data models.AdminUpdateGachaRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	if err := cnt.srv.UpdateGacha(gachaId, &data); err != nil {
		if err == models.ErrGachaNotFound {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
		} else if err == models.ErrGachaAlreadyExists {
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
		} else {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
		}
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Gacha updated successfully!"})
}

func (cnt *AdminController) GetAllGachas(ctx *gin.Context) {
	if gachas, err := cnt.srv.GetAllGachas(); err != nil {
		if err == models.ErrGachaNotFound {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": models.ErrGachaNotFound})
		} else {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
		}
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"GachaList": gachas})
	}
}

func (cnt *AdminController) GetGachaDetails(ctx *gin.Context) {
	gachaId := ctx.Param("gacha_id")
	if gachaId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "gacha_id is required"})
		ctx.Abort()
		return
	}

	if gacha, err := cnt.srv.FindGachaByID(gachaId); err != nil {
		if err == models.ErrGachaNotFound {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
		} else {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err})
			ctx.Abort()
		}
		return
	} else {
		ctx.JSON(http.StatusOK, gacha)
	}
}

// Market controllers ==============================================

func (cnt *AdminController) GetMarketHistory(ctx *gin.Context) {
	transactions, ok := cnt.srv.GetMarketHistory()
	if !ok {
		ctx.Status(http.StatusInternalServerError)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"MarketHistory": transactions})
}

func (cnt *AdminController) UpdateAuction(ctx *gin.Context) {
	ctx.Status(http.StatusNotImplemented)
}

func (cnt *AdminController) GetAllAuctions(ctx *gin.Context) {
	if auctions, ok := cnt.srv.GetAllAuctions(); ok {
		ctx.JSON(http.StatusOK, gin.H{"AuctionList": auctions})
		return
	}
	ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", models.ErrInternalServerError)
	ctx.Abort()
}

func (cnt *AdminController) GetAuctionDetails(ctx *gin.Context) {
	auctionId := ctx.Param("auction_id")
	if auctionId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "auction_id is required"})
		ctx.Abort()
		return
	}

	auction, exists := cnt.srv.FindAuctionByID(auctionId)
	if !exists {
		ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": "Auction not found"})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, auction)
}
