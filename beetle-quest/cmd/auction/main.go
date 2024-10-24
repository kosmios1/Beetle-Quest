package main

import (
	"github.com/gin-gonic/gin"

	"beetle-quest/internal/auction/controller"
	"beetle-quest/internal/auction/service"
)

func main() {
	r := gin.Default()

	var auctionController controller.AuctionController
	auctionController = controller.AuctionController{
		Service: service.AuctionService{
			UserRepo:    nil,
			GachaRepo:   nil,
			AuctionRepo: nil,
		},
	}

	r.POST("/auction", auctionController.CreateAuction)
	r.GET("/auction/:auction_id", auctionController.GetAuction)

	r.DELETE("/auction/:auction_id", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "NOT IMPLEMENTED"})
	})
	r.POST("/auction/:auction_id/bid", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "NOT IMPLEMENTED"})
	})

	r.Run()
}
