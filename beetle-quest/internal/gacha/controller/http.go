package controller

import (
	"beetle-quest/internal/gacha/service"
	"beetle-quest/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GachaController struct {
	service.GachaService
}

func NewGachaController(s service.GachaService) GachaController {
	return GachaController{
		GachaService: s,
	}
}

func (c *GachaController) Roll(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Not implemented yet!"})
}

func (c *GachaController) List(ctx *gin.Context) {
	gachas, ok := c.GachaService.GetAll()
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

	gacha, err := c.FindByID(id)
	if err != nil {
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
