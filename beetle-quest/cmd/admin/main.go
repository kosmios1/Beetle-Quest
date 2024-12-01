package main

import (
	"beetle-quest/pkg/httpserver"
	"beetle-quest/pkg/models"

	"github.com/gin-gonic/gin"

	entrypoint "beetle-quest/internal/admin/entrypoints"
	"beetle-quest/internal/admin/middleware"

	packageMiddleware "beetle-quest/pkg/middleware/secure"
)

func main() {
	httpserver.GenOwnCertAndKey("admin-service")

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(packageMiddleware.NewSecureMiddleware())

	r.LoadHTMLGlob("templates/*")

	cnt := entrypoint.NewAdminController()

	basePath := r.Group("/api/v1/admin")
	basePath.Use(middleware.CheckAdminJWTAuthorizationToken(models.AdminScope))
	{
		userPath := basePath.Group("/user")
		{
			userPath.GET("/get_all", cnt.GetAllUsers)
			userPath.GET("/:user_id", cnt.GetUserProfile)
			userPath.PATCH("/:user_id", cnt.UpdateUserProfile)
			userPath.GET("/:user_id/transaction_history", cnt.GetUserTransactionHistory)
			userPath.GET("/:user_id/auction/get_all", cnt.GetUserAuctionList)
		}

		gachaPath := basePath.Group("/gacha")
		{
			gachaPath.POST("/add", cnt.AddGacha)
			gachaPath.GET("/get_all", cnt.GetAllGachas)
			gachaPath.GET("/:gacha_id", cnt.GetGachaDetails)
			gachaPath.DELETE("/:gacha_id", cnt.DeleteGacha)
			gachaPath.PATCH("/:gacha_id", cnt.UpdateGacha)
		}

		marketPath := basePath.Group("/market")
		{
			marketPath.GET("/transaction_history", cnt.GetMarketHistory)
			auctionPath := marketPath.Group("/auction")
			{
				auctionPath.GET("/get_all", cnt.GetAllAuctions)
				auctionPath.GET("/:auction_id", cnt.GetAuctionDetails)
				auctionPath.PATCH("/:auction_id", cnt.UpdateAuction)
			}
		}
	}

	internalPath := r.Group("/api/v1/internal/admin")
	{
		internalPath.POST("/find_by_id", cnt.FindByID)
	}

	server := httpserver.SetupHTPPSServer(r)
	httpserver.ListenAndServe(server)
}
