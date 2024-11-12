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
	return &MarketController{srv: srv}
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
	// TODO: Steps to implement
	// 1. Get auction's id
	// 2. Get user_id by the gin context
	// 3. Check if user is the owner of the auction
	// 4. Check that the auction is not expired
	// 5. Check that no one bid to the auction
	// 6. Check that the auction is open less than x time
	// 7. Refund the bidders
	// 8. Delete the auction

}

func (c *MarketController) BidToAuction(ctx *gin.Context) {
	// TODO: Steps to implement
	// 1. Get auction's id
	// 2. Get user_id by the gin context
	// 3. Check if user is not the owner of the auction
	// 4. Check if user has enough bugscoin to bid
	// 5. Check if the auction is not expired
	// 6. Bid to the auction
}
