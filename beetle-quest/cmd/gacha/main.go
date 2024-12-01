package main

import (
	"beetle-quest/pkg/httpserver"
	middleware "beetle-quest/pkg/middleware/authorization"
	"beetle-quest/pkg/models"

	entrypoint "beetle-quest/internal/gacha/entrypoints"

	"github.com/gin-gonic/gin"

	packageMiddleware "beetle-quest/pkg/middleware/secure"
)

func main() {
	httpserver.GenOwnCertAndKey("gacha-service")

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(packageMiddleware.NewSecureMiddleware())

	r.LoadHTMLGlob("templates/*")

	cnt := entrypoint.NewGachaController()

	basePath := r.Group("/api/v1/gacha")
	basePath.Use(middleware.CheckJWTAuthorizationToken(models.GachaScope))
	{
		basePath.GET("/list", cnt.List)
		basePath.GET("/:gacha_id", cnt.GetGachaDetails)

		basePath.GET("/user/:user_id/list", cnt.GetUserGachaList)
		basePath.GET("/:gacha_id/user/:user_id", cnt.GetUserGachaDetails)
	}

	internalPath := r.Group("/api/v1/internal/gacha")
	{
		internalPath.POST("/create", cnt.CreateGacha)
		internalPath.POST("/update", cnt.UpdateGacha)
		internalPath.POST("/delete", cnt.DeleteGacha)

		internalPath.POST("/get_user_gachas", cnt.GetUserGachas)
		internalPath.POST("/remove_user_gachas", cnt.RemoveUserGachas)
		internalPath.POST("/add_gacha_to_user", cnt.AddGachaToUser)
		internalPath.POST("/remove_gacha_from_user", cnt.RemoveGachaFromUser)
		internalPath.POST("/find_by_id", cnt.FindByID)
		internalPath.GET("/get_all", cnt.GetAll)
	}

	server := httpserver.SetupHTPPSServer(r)
	httpserver.ListenAndServe(server)
}
