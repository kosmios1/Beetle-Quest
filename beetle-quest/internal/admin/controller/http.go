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

// User controllers =================================================

func (c *AdminController) GetAllUsers(ctx *gin.Context) {
	users, err := c.srv.GetAllUsers()
	if err != nil {
		// TODO: Use switch statement
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
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrUserNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		panic("unreachable code")
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *AdminController) UpdateUserProfile(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	if userId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	var data models.AdminUpdateUserAccount
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	if err := c.srv.UpdateUserProfile(userId, &data); err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrUserNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrUsernameOrEmailAlreadyExists:
			ctx.HTML(http.StatusConflict, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInvalidData:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		panic("unreachable code")
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "User profile updated successfully!"})
}

func (c *AdminController) GetUserTransactionHistory(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	if userId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	transactions, err := c.srv.GetUserTransactionHistory(userId)
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrTransactionNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		panic("code should not reach here")
	}

	// ctx.HTML(http.StatusOK, "userTransactionHistory.tmpl", gin.H{"TransactionList": transactions})
	ctx.JSON(http.StatusOK, gin.H{"TransactionList": transactions})
}

func (c *AdminController) GetUserAuctionList(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	if userId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	auctionList, err := c.srv.GetUserAuctionList(userId)
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrAuctionNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
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
		// TODO: Use switch statement
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
		// TODO: Use switch statement
		if err == models.ErrGachaNotFound {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
		} else {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
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
		// TODO: Use switch statement
		if err == models.ErrGachaNotFound {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
		} else if err == models.ErrGachaAlreadyExists {
			ctx.HTML(http.StatusConflict, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
		} else {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
		}
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Gacha updated successfully!"})
}

func (cnt *AdminController) GetAllGachas(ctx *gin.Context) {
	if gachas, err := cnt.srv.GetAllGachas(); err != nil {
		// TODO: Use switch statement
		if err == models.ErrGachaNotFound {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": models.ErrGachaNotFound})
		} else {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
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
		// TODO: Use switch statement
		if err == models.ErrGachaNotFound {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
		} else {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
		}
		return
	} else {
		ctx.JSON(http.StatusOK, gacha)
	}
}

// Market controllers ==============================================

func (cnt *AdminController) GetMarketHistory(ctx *gin.Context) {
	transactions, err := cnt.srv.GetMarketHistory()
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", err)
			ctx.Abort()
			return
		case models.ErrTransactionNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", err)
			ctx.Abort()
			return
		}
		panic("unreachable code")
	}
	ctx.JSON(http.StatusOK, gin.H{"MarketHistory": transactions})
}

func (cnt *AdminController) UpdateAuction(ctx *gin.Context) {
	ctx.Status(http.StatusNotImplemented)
}

func (cnt *AdminController) GetAllAuctions(ctx *gin.Context) {
	auctions, err := cnt.srv.GetAllAuctions()
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", err)
			ctx.Abort()
			return
		case models.ErrTransactionNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", err)
			ctx.Abort()
			return
		}
		panic("unreachable code")
	}

	ctx.JSON(http.StatusOK, gin.H{"AuctionList": auctions})
}

func (cnt *AdminController) GetAuctionDetails(ctx *gin.Context) {
	auctionId := ctx.Param("auction_id")
	if auctionId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "auction_id is required"})
		ctx.Abort()
		return
	}

	auction, err := cnt.srv.FindAuctionByID(auctionId)
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrAuctionNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		panic("unreachable code")
	}

	ctx.JSON(http.StatusOK, auction)
}

// Internal =================================================
func (c *AdminController) FindByID(ctx *gin.Context) {
	var data models.FindAdminByIDData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	admin, err := c.srv.FindByID(data.AdminID)
	if err != nil {
		// TODO: Use switch statement
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		} else if err == models.ErrAdminNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}
		panic("unreachable code")
	}
	ctx.JSON(http.StatusOK, admin)
}
