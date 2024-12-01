package main

import (
	"beetle-quest/pkg/httpserver"
	middleware "beetle-quest/pkg/middleware/authorization"
	"beetle-quest/pkg/models"
	"log"

	"github.com/gin-gonic/gin"

	entrypoint "beetle-quest/internal/market/entrypoints"

	packageMiddleware "beetle-quest/pkg/middleware/secure"
)

func main() {
	httpserver.GenOwnCertAndKey("market-service")

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(packageMiddleware.NewSecureMiddleware())

	r.LoadHTMLGlob("templates/*")

	cnt, err := entrypoint.NewMarketController()
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}

	basePath := r.Group("/api/v1/market")
	basePath.Use(middleware.CheckJWTAuthorizationToken(models.MarketScope))
	{
		basePath.POST("/bugscoin/buy", cnt.BuyBugscoin)
		basePath.GET("/gacha/roll", cnt.RollGacha)
		basePath.GET("/gacha/:gacha_id/buy", cnt.BuyGacha)

		auctionPath := basePath.Group("/auction")
		{
			auctionPath.POST("/", cnt.CreateAuction)
			auctionPath.GET("/list", cnt.AuctionList)
			auctionPath.GET("/:auction_id", cnt.AuctionDetail)
			auctionPath.POST("/:auction_id", cnt.AuctionDelete)
			auctionPath.POST("/:auction_id/bid", cnt.BidToAuction)
		}
	}

	internalPath := r.Group("/api/v1/internal/market")
	{

		internalPath.PATCH("/auction/update", cnt.UpdateAuction)
		internalPath.POST("/auction/find_by_id", cnt.FindAuctionByID)
		internalPath.GET("/auction/get_all", cnt.GetAllAuctions)
		internalPath.POST("/auction/get_user_auctions", cnt.GetUserAuctions)

		internalPath.GET("/get_transaction_history", cnt.GetTransactionHistory)
		internalPath.POST("/get_user_transaction_history", cnt.GetUserTransactionHistory)
		internalPath.POST("/delete_user_transaction_history", cnt.DeleteUserTransactionHistory)
	}

	server := httpserver.SetupHTPPSServer(r)
	httpserver.ListenAndServe(server)
}
