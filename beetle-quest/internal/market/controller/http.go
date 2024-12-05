package controller

import (
	service "beetle-quest/internal/market/service"
	"beetle-quest/pkg/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MarketController struct {
	srv *service.MarketService
}

func NewMarketController(srv *service.MarketService) *MarketController {
	return &MarketController{
		srv: srv,
	}
}

func (c *MarketController) BuyBugscoin(ctx *gin.Context) {
	var buyBugscoinRequest models.BuyBugscoinRequest
	if err := ctx.ShouldBindJSON(&buyBugscoinRequest); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	amount, err := strconv.Atoi(buyBugscoinRequest.Amount)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidUserID})
		ctx.Abort()
		return
	}

	if err := c.srv.AddBugsCoin(userId.(string), int64(amount)); err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrUserNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrMaxMoneyExceeded:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}
	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Bugscoin added successfully"})
}

func (c *MarketController) RollGacha(ctx *gin.Context) {
	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrUserNotFound})
		ctx.Abort()
		return
	}

	gacha, msg, err := c.srv.RollGacha(userId.(string))
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrUserNotFound, models.ErrGachaNotFound, models.ErrUserAlreadyHasGacha:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrNotEnoughMoneyToRollGacha:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	if gacha == nil {
		ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": msg})
	} else {
		ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": msg, "HiddenData": gacha.GachaID.String(), "HiddenData2": gacha.Rarity.String()})
	}
}

func (c *MarketController) BuyGacha(ctx *gin.Context) {
	gachaId := ctx.Param("gacha_id")
	if gachaId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidGachaID})
		ctx.Abort()
		return
	}

	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidUserID})
		ctx.Abort()
		return
	}

	gacha, err := c.srv.BuyGacha(userId.(string), gachaId)
	if err != nil {
		switch err {
		case models.ErrGachaNotFound, models.ErrUserNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrUserAlreadyHasGacha, models.ErrNotEnoughMoneyToBuyGacha:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Gacha bought successfully", "HiddenData": gacha.GachaID.String()})
}

func (c *MarketController) CreateAuction(ctx *gin.Context) {
	var data models.CreateAuctionRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	const layout = "2006-01-02T15:04"
	endTime, err := time.Parse(layout, data.EndTime)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidTimeFormat})
		ctx.Abort()
		return
	}

	uid, ok := ctx.Get("user_id")
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidUserID})
		ctx.Abort()
		return
	}

	auction, err := c.srv.CreateAuction(uid.(string), data.GachaID, endTime)
	if err != nil {
		switch err {
		case models.ErrGachaNotFound, models.ErrUserNotFound, models.ErrAuctionNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInvalidData, models.ErrUserDoesNotOwnGacha, models.ErrInvalidEndTime:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrAuctionAltreadyExists, models.ErrGachaAlreadyAuctioned:
			ctx.HTML(http.StatusConflict, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Auction created successfully", "HiddenData": auction.AuctionID.String()})
}

func (c *MarketController) AuctionList(ctx *gin.Context) {
	auctions, err := c.srv.RetrieveAuctionTemplateList()
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
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	ctx.HTML(http.StatusOK, "market.tmpl", gin.H{"Auctions": auctions})
}

func (c *MarketController) AuctionDetail(ctx *gin.Context) {
	auctionId := ctx.Param("auction_id")
	if auctionId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	auction, bids, err := c.srv.GetAuctionDetails(auctionId)
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
			ctx.Abort()
			return
		case models.ErrBidsNotFound, models.ErrAuctionNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	data := struct {
		Auction *models.Auction
		Bids    []models.Bid
	}{
		Auction: auction,
		Bids:    bids,
	}
	ctx.HTML(http.StatusOK, "auctionDetails.tmpl", data)
}

func (c *MarketController) AuctionDelete(ctx *gin.Context) {
	var data struct {
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&data); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidPassword})
		ctx.Abort()
		return
	}

	aid := ctx.Param("auction_id")
	if aid == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidAuctionID})
		ctx.Abort()
		return
	}

	uid, ok := ctx.Get("user_id")
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidUserID})
		ctx.Abort()
	}

	if err := c.srv.DeleteAuction(uid.(string), aid, data.Password); err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrAuctionNotFound, models.ErrUserNotFound, models.ErrBidsNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrUserNotOwnerOfAuction, models.ErrInvalidPassword, models.ErrAuctionEnded, models.ErrAuctionIsTooCloseToEnd, models.ErrAuctionHasBids:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Auction deleted successfully"})
}

func (c *MarketController) BidToAuction(ctx *gin.Context) {
	var data models.BidRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	aid := ctx.Param("auction_id")
	if aid == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidAuctionID})
		ctx.Abort()
		return
	}

	uid, ok := ctx.Get("user_id")
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidUserID})
		ctx.Abort()
		return
	}

	bidAmount, err := strconv.Atoi(data.BidAmount)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidBidAmount})
		ctx.Abort()
		return
	}

	err = c.srv.MakeBid(uid.(string), aid, int64(bidAmount))
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrAuctionNotFound, models.ErrUserNotFound:
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrOwnerCannotBid, models.ErrBidAmountNotEnough, models.ErrNotEnoughMoneyToBid, models.ErrAuctionEnded:
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err, err: %s", err.Error())
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Bid successfully"})
}

// Internal ==========================================================================================================

func (c *MarketController) GetAllAuctions(ctx *gin.Context) {
	auctions, err := c.srv.GetAuctionList()
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrAuctionNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}
	ctx.JSON(http.StatusOK, models.GetAllAuctionDataResponse{AuctionList: auctions})
}

func (c *MarketController) UpdateAuction(ctx *gin.Context) {
	var data models.Auction
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	if err := c.srv.UpdateAuction(&data); err != nil {
		switch err {
		case models.ErrAuctionNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInternalServerError:
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInvalidData:
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}
}

func (c *MarketController) GetTransactionHistory(ctx *gin.Context) {
	transactions, err := c.srv.GetAllTransactions()
	if err != nil {
		switch err {
		case models.ErrTransactionNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInternalServerError:
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	ctx.JSON(http.StatusOK, models.GetAllTransactionDataResponse{TransactionHistory: transactions})
}

func (c *MarketController) GetUserAuctions(ctx *gin.Context) {
	var data models.GetUserAuctionsData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	auctions, err := c.srv.GetAuctionListOfUser(data.UserID)
	if err != nil {
		switch err {
		case models.ErrAuctionNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInternalServerError:
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	ctx.JSON(http.StatusOK, models.GetUserAuctionsDataResponse{AuctionList: auctions})
}

func (c *MarketController) FindAuctionByID(ctx *gin.Context) {
	var data models.FindAuctionByIDData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	auction, err := c.srv.FindByID(data.AuctionID.String())
	if err != nil {
		switch err {
		case models.ErrAuctionNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrInternalServerError:
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}
	ctx.JSON(http.StatusOK, models.FindAuctionByIDDataResponse{Auction: auction})
}

func (c *MarketController) GetUserTransactionHistory(ctx *gin.Context) {
	var data models.GetUserTransactionHistoryData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	auctions, err := c.srv.GetUserTransactionHistory(data.UserID)
	if err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrTransactionNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}
	ctx.JSON(http.StatusOK, models.GetUserTransactionHistoryDataResponse{TransactionHistory: auctions})
}

func (c *MarketController) DeleteUserTransactionHistory(ctx *gin.Context) {
	var data models.DeleteUserTransactionHistoryData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	if err := c.srv.DeleteUserTransactionHistory(data.UserID); err != nil {
		switch err {
		case models.ErrInternalServerError:
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		case models.ErrTransactionNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		log.Panicf("Unreachable code, err: %s", err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "Transaction history deleted successfully"})
}
