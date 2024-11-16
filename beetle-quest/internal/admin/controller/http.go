package controller

import (
	"beetle-quest/internal/admin/service"
	"beetle-quest/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	srv *service.AdminService
}

func NewAdminController(srv *service.AdminService) *AdminController {
	return &AdminController{
		srv: srv,
	}
}

func (c *AdminController) FindByID(ctx *gin.Context) {
	var data models.FindAdminByIDData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	admin, exists := c.srv.FindByID(data.AdminID)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}
	ctx.JSON(http.StatusOK, admin)
}
