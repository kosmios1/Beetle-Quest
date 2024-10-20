package controller

import (
	"encoding/base64"
	"gacha-app/internal/auction/service"
	"gacha-app/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuctionController struct {
	Service service.AuctionService
}

func (c *AuctionController) CreateAuction(ctx *gin.Context) {
	var req models.CreateAuctionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endTime, err := time.ParseInLocation(time.RFC3339, req.EndTime, time.UTC)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.ErrCouldNotParseTime})
		return
	}

	ownerID, err := base64.StdEncoding.DecodeString(req.OwnerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.ErrCouldNotDecodeUserID})
		return
	}

	gachaID, err := base64.StdEncoding.DecodeString(req.GachaID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.ErrCouldNotDecodeGachaID})
		return
	}

	auction, err := c.Service.CreateAuction(ownerID, gachaID, endTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, models.CreateAuctionRes{Auction: auction})
}
