package main

import (
	"beetle-quest/internal/market/controller"
	"beetle-quest/internal/market/service"
	"beetle-quest/pkg/middleware"
	"beetle-quest/pkg/utils"
	"log"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"

	arepo "beetle-quest/internal/market/repository"
	grepo "beetle-quest/pkg/repositories/serviceHttp/gacha"
	urepo "beetle-quest/pkg/repositories/serviceHttp/user"
)

func main() {
	utils.GenOwnCertAndKey("market")

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(secure.New(secure.Config{
		SSLRedirect:           true,
		IsDevelopment:         false,
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		AllowedHosts:          []string{},
	}))

	r.LoadHTMLGlob("templates/*")

	cnt := controller.NewMarketController(
		service.NewMarketService(
			urepo.NewUserRepo(),
			grepo.NewGachaRepo(),
			arepo.NewMarketRepo(),
		),
	)

	basePath := r.Group("/api/v1/market")
	basePath.Use(middleware.CheckJWTAuthorizationToken())
	{
		basePath.GET("/", cnt.Market)
		basePath.POST("/bugscoin/buy", cnt.BuyBugscoin)
		basePath.GET("/gacha/roll", cnt.RollGacha)
		basePath.GET("/gacha/:gacha_id/buy", cnt.BuyGacha)

		auctionPath := basePath.Group("/auction")
		{
			auctionPath.POST("/", cnt.CreateAuction)
			auctionPath.GET("/list", cnt.AuctionList)
			auctionPath.GET("/:auction_id", cnt.AuctionDetail)
			auctionPath.DELETE("/:auction_id", cnt.AuctionDelete)
			auctionPath.POST("/:auction_id/bid", cnt.BidToAuction)
		}
	}

	internalPath := r.Group("/api/v1/internal/market")
	{
		internalPath.POST("/get_user_transaction_history", cnt.GetUserTransactionHistory)
	}

	server := utils.SetupHTPPSServer(r)
	if err := server.ListenAndServeTLS("/serverCert.pem", "/serverKey.pem"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
