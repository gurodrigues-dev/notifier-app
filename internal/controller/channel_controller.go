package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
	"github.com/gurodrigues-dev/notifier-app/internal/usecase"
	"github.com/jinzhu/gorm"
)

type ChannelController struct {
	logger contracts.Logger
}

func NewChannelController(
	logger contracts.Logger,
) *ChannelController {
	return &ChannelController{
		logger: logger,
	}
}

func (cc *ChannelController) CreateChannel(httpContext *gin.Context) {
	var requestParams entity.Channel
	if err := httpContext.BindJSON(&requestParams); err != nil {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid body content"})
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestParams); err != nil {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecase := usecase.NewCreateChannelUsecase(
		infra.App.Repositories.ChannelRepository,
		infra.App.Email,
		cc.logger,
	)

	channel, err := usecase.CreateChannel(&requestParams)
	if err != nil {
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpContext.JSON(http.StatusCreated, channel)
	return
}

func (cc *ChannelController) FindById(httpContext *gin.Context) {
	id := httpContext.Param("id")
	if id == "" {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": "target is required"})
		return
	}

	usecase := usecase.NewGetChannelByIDUsecase(
		infra.App.Repositories.ChannelRepository,
		cc.logger,
	)

	channel, err := usecase.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httpContext.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpContext.JSON(http.StatusOK, channel)
	return
}

func (cc *ChannelController) FindByGroup(httpContext *gin.Context) {
	groupID := httpContext.Param("group")
	if groupID == "" {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": "group is required"})
		return
	}

	usecase := usecase.NewListChannelsByGroupUsecase(
		infra.App.Repositories.ChannelRepository,
		cc.logger,
	)

	channel, err := usecase.ListByGroupID(groupID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httpContext.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpContext.JSON(http.StatusOK, channel)
	return
}

func (cc *ChannelController) FindByPlatform(httpContext *gin.Context) {
	platformID := httpContext.Param("platform")
	if platformID == "" {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": "group is required"})
		return
	}

	usecase := usecase.NewListChannelsByPlatformUsecase(
		infra.App.Repositories.ChannelRepository,
		cc.logger,
	)

	channel, err := usecase.ListByPlatformID(platformID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httpContext.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpContext.JSON(http.StatusOK, channel)
	return
}

func (cc *ChannelController) DeleteById(httpContext *gin.Context) {
	id := httpContext.Param("id")
	if id == "" {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": "target is required"})
		return
	}

	usecase := usecase.NewDeleteChannelByIDUsecase(
		infra.App.Repositories.ChannelRepository,
		cc.logger,
	)

	err := usecase.DeleteByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httpContext.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpContext.JSON(http.StatusOK, "channel deleted successfully")
	return
}
