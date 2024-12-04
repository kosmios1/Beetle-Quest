package main

import (
	"beetle-quest/pkg/httpserver"
	middleware "beetle-quest/pkg/middleware/authorization"
	"beetle-quest/pkg/models"

	entrypoint "beetle-quest/internal/user/entrypoints"
	internalMiddleware "beetle-quest/internal/user/middleware"

	"github.com/gin-gonic/gin"

	packageMiddleware "beetle-quest/pkg/middleware/secure"
)

func main() {
	httpserver.GenOwnCertAndKey("user-service")

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(packageMiddleware.NewSecureMiddleware())

	r.LoadHTMLGlob("templates/*")

	cnt := entrypoint.NewUserController()

	r.GET("/userinfo", middleware.CheckJWTAuthorizationToken(models.UserScope), cnt.UserInfo)

	basePath := r.Group("/api/v1/user")
	basePath.Use(middleware.CheckJWTAuthorizationToken(models.UserScope))
	{
		accountGroup := basePath.Group("/account")
		accountGroup.Use(internalMiddleware.CheckUserID())
		{
			accountGroup.GET("/:user_id", cnt.GetUserAccountDetails)
			accountGroup.PATCH("/:user_id", cnt.UpdateUserAccountDetails)
			accountGroup.POST("/:user_id", cnt.DeleteUserAccount)
		}
	}

	internalPath := r.Group("/api/v1/internal/user")
	{
		internalPath.GET("/get_all", cnt.GetAllUsers)
		internalPath.POST("/create", cnt.CreateUser)
		internalPath.POST("/update", cnt.UpdateUser)

		internalPath.POST("/find_by_id", cnt.FindByID)
		internalPath.POST("/find_by_username", cnt.FindByUsername)
	}

	server := httpserver.SetupHTPPSServer(r)
	httpserver.ListenAndServe(server)
}
