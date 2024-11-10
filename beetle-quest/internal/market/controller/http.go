package controller

import (
	service "beetle-quest/internal/market/service"
	"beetle-quest/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MarketController struct {
	service.MarketService
}

func (c *MarketController) BuyBugscoin(ctx *gin.Context) {
	var buyBugscoinRequest models.BuyBugscoinRequest
	if err := ctx.ShouldBindJSON(buyBugscoinRequest); err != nil {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"error": err.Error()})
		return
	}

	// TODO: Steps to implement
	// -. Receive amount of bugscoin to add
	// 2. Get user_id by the gin context
	// 3. Call userRepo to add bugscoin to the user
}

func (c *MarketController) BuyGacha(ctx *gin.Context) {
	// TODO: Steps to implement
	// 1. Receive the gacha's id to buy
	// 2. Get user_id by the gin context
	// 3. Check gacha's price - user's bugscoin > 0
	// 4. Subtract gacha's price from user's bugscoin
	// 5. Add gacha to user's inventory
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
