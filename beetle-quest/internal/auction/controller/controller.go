package controller

import (
	"beetle-quest/internal/auction/service"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuctionController struct {
	Service service.AuctionService
}

func (c *AuctionController) CreateAuction(ctx *gin.Context) {
	var req models.CreateAuctionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endTime, err := time.ParseInLocation(time.RFC3339, req.EndTime, time.UTC)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.ErrCouldNotParseTime})
		return
	}

	ownerUUID, err := utils.Parse(ctx.Param("user_uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	gachaUUID, err := utils.Parse(req.GachaUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	auction, err := c.Service.CreateAuction(ownerUUID, gachaUUID, endTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, models.CreateAuctionResponse{Auction: auction})
}

func (c *AuctionController) GetAuction(ctx *gin.Context) {
	auctionUUID, err := utils.Parse(ctx.Param("auction_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	auction, err := c.Service.GetAuction(auctionUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.GetAuctionResponse{Auction: auction})
}
