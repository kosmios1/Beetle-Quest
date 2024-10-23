package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"beetle-quest/internal/auction/controller"
	"beetle-quest/internal/auction/service"
	mock_repositories "beetle-quest/mock/repositories"
)

func main() {
	r := gin.Default()

	var auctionController controller.AuctionController
	if _, ok := os.LookupEnv("DEBUG"); ok {
		auctionController = controller.AuctionController{
			Service: service.AuctionService{
				UserRepo:    mock_repositories.MockUserRepo{},
				GachaRepo:   mock_repositories.MockGachaRepo{},
				AuctionRepo: mock_repositories.MockAuctionRepo{},
			},
		}
	} else {
		panic("Not yet implemented")
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
