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
	gachas, err := c.srv.GetAll()
	if err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
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
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidGachaID})
		ctx.Abort()
		return
	}

	gacha, err := c.srv.FindByID(id)
	if err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
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

func (cnt *GachaController) GetUserGachaList(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	if userId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidUserID})
		ctx.Abort()
		return
	}

	gacha, err := cnt.srv.GetUserGachasStr(userId)
	if err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		} else if err == models.ErrInvalidData {
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		} else {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gacha)
}

func (cnt *GachaController) GetUserGachaDetails(ctx *gin.Context) {
	gachaId := ctx.Param("gacha_id")
	if gachaId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidGachaID})
		ctx.Abort()
		return
	}

	userId := ctx.Param("user_id")
	if userId == "" {
		ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": models.ErrInvalidUserID})
		ctx.Abort()
		return
	}

	gacha, err := cnt.srv.GetUserGachaDetails(userId, gachaId)
	if err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.HTML(http.StatusInternalServerError, "errorMsg.tmpl", gin.H{"Error": err.Error()})
			ctx.Abort()
			return
		} else if err == models.ErrInvalidData {
			ctx.HTML(http.StatusBadRequest, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		} else {
			ctx.HTML(http.StatusNotFound, "errorMsg.tmpl", gin.H{"Error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gacha)
}

// Internal API ============================================================================================================

func (c *GachaController) CreateGacha(ctx *gin.Context) {
	var data models.Gacha
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		return
	}

	if err := c.srv.CreateGacha(&data); err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		} else if err == models.ErrGachaAlreadyExists {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"Error": err.Error()})
			return
		}
		panic("Unreachable code")
	}

	ctx.Status(http.StatusOK)
}

func (c *GachaController) UpdateGacha(ctx *gin.Context) {
	var data models.Gacha
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		return
	}

	if err := c.srv.UpdateGacha(&data); err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		} else if err == models.ErrGachaNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		} else if err == models.ErrGachaAlreadyExists {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"Error": err.Error()})
			return
		}
		panic("Unreachable code")
	}

	ctx.Status(http.StatusOK)
}

func (c *GachaController) DeleteGacha(ctx *gin.Context) {
	var data models.Gacha
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		return
	}

	if err := c.srv.DeleteGacha(&data); err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		} else if err == models.ErrGachaNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}
		panic("Unreachable code")
	}

	ctx.Status(http.StatusOK)
}

func (c *GachaController) GetAll(ctx *gin.Context) {
	gachas, err := c.srv.GetAll()
	if err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		} else if err == models.ErrGachaNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}
		panic("Unreachable code")
	}

	ctx.JSON(http.StatusOK, models.GetAllGachasDataResponse{GachaList: gachas})
}

func (c *GachaController) GetUserGachas(ctx *gin.Context) {
	var data models.GetUserGachasData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		return
	}

	gachas, err := c.srv.GetUserGachas(data.UserID)
	if err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		} else if err == models.ErrUserNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}
		panic("Unreachable code")
	}

	ctx.JSON(http.StatusOK, models.GetUserGachasDataResponse{GachaList: gachas})
}

func (c *GachaController) AddGachaToUser(ctx *gin.Context) {
	var data models.AddGachaToUserData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		return
	}

	if err := c.srv.AddGachaToUser(data.UserID, data.GachaID); err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		} else if err == models.ErrUserAlreadyHasGacha {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"Error": err.Error()})
			return
		}
		panic("Unreachable code")
	}

	ctx.Status(http.StatusOK)
}

func (c *GachaController) RemoveGachaFromUser(ctx *gin.Context) {
	var data models.RemoveGachaFromUserData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		return
	}

	if err := c.srv.RemoveGachaFromUser(data.UserID, data.GachaID); err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		} else if err == models.ErrRetalationGachaUserNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}
		panic("Unreachable code")
	}

	ctx.Status(http.StatusOK)
}

func (c *GachaController) FindByID(ctx *gin.Context) {
	var req models.FindGachaByIDData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	gacha, err := c.srv.FindByID(req.GachaID)
	if err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		} else if err == models.ErrGachaNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}
		panic("Unreachable code")
	}

	ctx.JSON(http.StatusOK, gacha)
}

func (c *GachaController) RemoveUserGachas(ctx *gin.Context) {
	var req models.RemoveUserGachasData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": models.ErrInvalidData})
		ctx.Abort()
		return
	}

	if err := c.srv.RemoveUserGachas(req.UserID); err != nil {
		// TODO: use switch
		if err == models.ErrInternalServerError {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		} else if err == models.ErrRetalationGachaUserNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}
		panic("Unreachable code")
	}

	ctx.Status(http.StatusOK)
}
