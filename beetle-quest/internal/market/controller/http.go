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
}

func (c *MarketController) BuyGacha(ctx *gin.Context) {
}

func (c *MarketController) CreateAuction(ctx *gin.Context) {
}

func (c *MarketController) AuctionList(ctx *gin.Context) {
}

func (c *MarketController) AuctionDetail(ctx *gin.Context) {
}

func (c *MarketController) AuctionDelete(ctx *gin.Context) {
}

func (c *MarketController) BidTOAuction(ctx *gin.Context) {
}
