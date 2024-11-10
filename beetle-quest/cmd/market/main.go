package main

import (
	"beetle-quest/internal/market/controller"
	"beetle-quest/internal/market/service"
	"beetle-quest/pkg/middleware"
	repository "beetle-quest/pkg/repositories/httpImpl"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	// TODO: Uncomment this when having a valid SSL certificate
	// r.Use(secure.New(secure.Config{
	// 	SSLRedirect:           true,
	// 	IsDevelopment:         false,
	// 	STSSeconds:            315360000,
	// 	STSIncludeSubdomains:  true,
	// 	FrameDeny:             true,
	// 	ContentTypeNosniff:    true,
	// 	BrowserXssFilter:      true,
	// 	ContentSecurityPolicy: "default-src 'self'",
	// 	IENoOpen:              true,
	// 	SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	// 	AllowedHosts:          []string{},
	// }))

	r.LoadHTMLGlob("templates/*")

	cnt := controller.MarketController{
		MarketService: service.MarketService{
			UserRepo:    repository.NewUserRepo(),
			GachaRepo:   repository.NewGachaRepo(),
			AuctionRepo: repository.NewAuctionRepo(),
		},
	}

	basePath := r.Group("/api/v1/market")
	basePath.Use(middleware.CheckJWTAuthorizationToken())
	{
		basePath.POST("/bugscoin/buy", cnt.BuyBugscoin)
		basePath.GET("/gacha/:gacha_id/buy", nil)

		auctionPath := r.Group("/auction")
		{
			auctionPath.POST("/", cnt.CreateAuction)
			auctionPath.GET("/list", cnt.AuctionList)
			auctionPath.GET("/:auction_id", cnt.AuctionDetail)
			auctionPath.DELETE("/:auction_id", cnt.AuctionDelete)
			auctionPath.POST("/:auction_id/bid", cnt.BidToAuction)
		}
	}

	r.Run()
}
