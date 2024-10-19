package main

import (
	"github.com/gin-gonic/gin"

	"gacha-app/internal/auction/controller"
	"gacha-app/internal/auction/service"
)

func main() {
	r := gin.Default()

	auctionController := controller.AuctionController{
		Service: service.AuctionService{
			UserRepo:    nil,
			GachaRepo:   nil,
			AuctionRepo: nil,
		},
	}

	r.POST("/auction", auctionController.CreateAuction)

	// TODO: Implement
	r.GET("/auction/:auction_id", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "NOT IMPLEMENTED"})
	})
	r.DELETE("/auction/:auction_id", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "NOT IMPLEMENTED"})
	})
	r.POST("/auction/:auction_id/bid", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "NOT IMPLEMENTED"})
	})

	r.Run()
}
