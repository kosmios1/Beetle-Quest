package main

import (
	"beetle-quest/pkg/utils"
	"log"
	"net/http"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"

	"beetle-quest/internal/admin/controller"
	"beetle-quest/internal/admin/middleware"
	"beetle-quest/internal/admin/repository"
	"beetle-quest/internal/admin/service"

	grepo "beetle-quest/pkg/repositories/serviceHttp/gacha"
	mrepo "beetle-quest/pkg/repositories/serviceHttp/market"
	urepo "beetle-quest/pkg/repositories/serviceHttp/user"
)

func main() {
	utils.GenOwnCertAndKey("admin")

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

	cnt := controller.NewAdminController(
		service.NewAdminService(
			repository.NewAdminRepo(),
			mrepo.NewMarketRepo(),
			urepo.NewUserRepo(),
			grepo.NewGachaRepo(),
		),
	)

	basePath := r.Group("/api/v1/admin")
	basePath.Use(middleware.CheckAdminJWTAuthorizationToken())
	{
		userPath := basePath.Group("/user")
		{
			userPath.GET("/get_all", cnt.GetAllUsers)
			userPath.GET("/:user_id", cnt.GetUserProfile)
			userPath.PATCH("/:user_id", cnt.UpdateUserProfile)
			userPath.GET("/:user_id/transaction_history", cnt.GetUserTransactionHistory)
			userPath.GET("/:user_id/market_history", func(ctx *gin.Context) {
				// TODO: Don't know what this means, it's inside the user stories
				ctx.Status(http.StatusNotImplemented)
			})
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
			marketPath.GET("/history", cnt.GetMarketHistory)
			auctionPath := marketPath.Group("/auction")
			{
				auctionPath.GET("/get_all", cnt.GetAllAuctions)
				auctionPath.GET("/:auction_id", cnt.GetAuctionDetails)
				auctionPath.PATCH("/:auction_id", cnt.UpdateAuction) // TODO: Implement this
			}
		}
	}

	internalPath := r.Group("/api/v1/internal/admin")
	{
		internalPath.POST("/find_by_id", cnt.FindByID)
	}

	server := utils.SetupHTPPSServer(r)
	if err := server.ListenAndServeTLS("/serverCert.pem", "/serverKey.pem"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
