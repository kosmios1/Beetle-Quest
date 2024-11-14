package controller

import (
	"beetle-quest/internal/gacha/service"
	"beetle-quest/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GachaController struct {
	srv *service.GachaService
}

func NewGachaController(s *service.GachaService) *GachaController {
	return &GachaController{
		srv: s,
	}
}

func (c *GachaController) List(ctx *gin.Context) {
	gachas, ok := c.srv.GetAll()
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": models.ErrGachaNotFound})
		return
	}

	gachaList := []models.GetGachaDetailsResponse{}
	for _, gacha := range gachas {
		gachaList = append(gachaList, models.GetGachaDetailsResponse{
			GachaID:   gacha.GachaID.String(),
			Name:      gacha.Name,
			Rarity:    gacha.Rarity.String(),
			Price:     gacha.Price,
			ImagePath: gacha.ImagePath,
		})
	}

	ctx.HTML(http.StatusOK, "gachaList.tmpl", gin.H{"GachaList": gachaList})
}

func (c *GachaController) GetGachaDetails(ctx *gin.Context) {
	id := ctx.Param("gacha_id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidGachaID})
		return
	}

	gacha, ok := c.srv.FindByID(id)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": models.ErrGachaNotFound})
		return
	}

	ctx.JSON(http.StatusOK, models.GetGachaDetailsResponse{
		GachaID:   gacha.GachaID.String(),
		Name:      gacha.Name,
		Rarity:    string(gacha.Rarity),
		Price:     gacha.Price,
		ImagePath: gacha.ImagePath,
	})
}

// Internal API ============================================================================================================

func (c *GachaController) GetAll(ctx *gin.Context) {
	gachas, ok := c.srv.GetAll()
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": models.ErrGachaNotFound})
		return
	}

	ctx.JSON(http.StatusOK, models.GetAllGachasDataResponse{GachaList: gachas})
}

func (c *GachaController) GetUserGachas(ctx *gin.Context) {
	var data models.GetUserGachasData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Invalid data submitted!"})
		return
	}

	gachas, ok := c.srv.GetUserGachas(data.UserID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": models.ErrGachaNotFound})
		return
	}

	ctx.JSON(http.StatusOK, models.GetUserGachasDataResponse{GachaList: gachas})
}

func (c *GachaController) AddGachaToUser(ctx *gin.Context) {
	var data models.AddGachaToUserData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Invalid data submitted!"})
		return
	}

	if err := c.srv.AddGachaToUser(data.UserID, data.GachaID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *GachaController) FindByID(ctx *gin.Context) {
	var req models.FindGachaByIDData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Wrong inputs passed to the request!"})
		ctx.Abort()
		return
	}

	gacha, ok := c.srv.FindByID(req.GachaID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": models.ErrGachaNotFound})
		return
	}

	ctx.JSON(http.StatusOK, gacha)
}
