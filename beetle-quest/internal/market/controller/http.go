package controller

import (
	service "beetle-quest/internal/market/service"
	"beetle-quest/pkg/models"
	"net/http"

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

	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": "user_id not correct!"})
		ctx.Abort()
		return
	}

	if err := c.srv.AddBugsCoin(userId.(string), buyBugscoinRequest.Amount); err != nil {
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
	// TODO: Steps to implement
	// 1. Receive the gacha's id to auction
	// 2. Get user_id by the gin context
	// 3. Check if user own gacha
	// 4. Create auction
	// 5. Add auction to the system
}

func (c *MarketController) AuctionList(ctx *gin.Context) {
	// TODO: Steps to implement
	// 1. Get all auctions
	// 2. Return all auctions
}

func (c *MarketController) AuctionDetail(ctx *gin.Context) {
	// TODO: Steps to implement
	// 1. Get auction's id
	// 2. Get auction's detail
	// 3. Return auction's detail
}

func (c *MarketController) AuctionDelete(ctx *gin.Context) {
	// TODO: Steps to implement
	// 1. Get auction's id
	// 2. Get user_id by the gin context
	// 3. Check if user is the owner of the auction
	// 4. Check that the auction is not expired
	// 5. Check that no one bid to the auction
	// 6. Check that the auction is open less than x time
	// 7. Delete the auction
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
