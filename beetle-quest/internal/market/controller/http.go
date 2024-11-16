package controller

import (
	service "beetle-quest/internal/market/service"
	"beetle-quest/pkg/models"
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

func (c *MarketController) Market(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "market.tmpl", gin.H{"userID": ctx.MustGet("userID")})
}

func (c *MarketController) BuyBugscoin(ctx *gin.Context) {
	var buyBugscoinRequest models.BuyBugscoinRequest
	if err := ctx.ShouldBindJSON(&buyBugscoinRequest); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "wrong request format!"})
		ctx.Abort()
		return
	}

	amount, err := strconv.Atoi(buyBugscoinRequest.Amount)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "amount not correct!"})
		ctx.Abort()
		return
	}

	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "user_id not correct!"})
		ctx.Abort()
		return
	}

	if err := c.srv.AddBugsCoin(userId.(string), int64(amount)); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
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

	msg, err := c.srv.RollGacha(userId.(string))
	if err != nil {
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": msg})
}

func (c *MarketController) BuyGacha(ctx *gin.Context) {
	gachaId := ctx.Param("gacha_id")
	if gachaId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "gacha_id not correct!"})
		ctx.Abort()
		return
	}

	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "user_id not correct!"})
		ctx.Abort()
		return
	}

	if err := c.srv.BuyGacha(userId.(string), gachaId); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Gacha bought successfully"})
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

	if err := c.srv.CreateAuction(uid.(string), data.GachaID, endTime); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Auction created successfully"})
}

func (c *MarketController) AuctionList(ctx *gin.Context) {
	auctions, err := c.srv.RetrieveAuctionTemplateList()
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
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

	auction, exists := c.srv.FindByID(auctionId)
	if !exists {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrAuctionNotFound})
		ctx.Abort()
		return
	}

	bids, ok := c.srv.GetBidListOfAuctionID(auctionId)
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrBidsNotFound})
		ctx.Abort()
		return
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
	password, ok := ctx.GetQuery("password")
	if !ok {
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

	if err := c.srv.DeleteAuction(uid.(string), aid, password); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
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
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "successMsg.tmpl", gin.H{"Message": "Bid successfully"})
}

// Internal ==========================================================================================================

func (c *MarketController) GetUserTransactionHistory(ctx *gin.Context) {
	var data models.GetUserTransactionHistoryData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	auctions, ok := c.srv.GetUserTransactionHistory(data.UserID)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": models.ErrTransactionNotFound})
		ctx.Abort()
		return
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

	if ok := c.srv.DeleteUserTransactionHistory(data.UserID); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": models.ErrCouldNotDelete})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "Transaction history deleted successfully"})
}
