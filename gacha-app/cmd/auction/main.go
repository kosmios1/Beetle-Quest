package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"gacha-app/internal/auction/controller"
	"gacha-app/internal/auction/service"
	mock_repositories "gacha-app/mock/repositories"
)

func main() {
	r := gin.Default()

	var auctionController controller.AuctionController
	if _, ok := os.LookupEnv("DEBUG"); ok {
		auctionController = controller.AuctionController{
			Service: service.AuctionService{
				UserRepo:    mock_repositories.NewMockUserRepo(),
				GachaRepo:   mock_repositories.NewMockGachaRepo(),
				AuctionRepo: mock_repositories.NewMockAuctionRepo(),
			},
		}
	} else {
		panic("Not yet implemented")
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
