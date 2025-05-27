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

// CreateChannel godoc
// @Summary Create a new channel
// @Description Creates a new channel based on the provided channel data.
// @Tags channel
// @Accept json
// @Produce json
// @Param channel body entity.Channel true "Channel request body"
// @Success 201 {object} entity.Channel "Channel created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /channel [post]
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

// FindById godoc
// @Summary Get channel by ID
// @Description Retrieves a channel by its unique identifier.
// @Tabs channel
// @Produce json
// @Param id path string true "Channel ID"
// @Success 200 {object} entity.Channel "Channel retrieved successfully"
// @Failure 400 {object} map[string]string "Channel ID is required"
// @Failure 404 {object} map[string]string "Channel not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /channel/{id} [get]
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

// FindByGroup godoc
// @Summary Get channels by group
// @Description Retrieves a list of channels associated with a specific group ID.
// @Tags channel
// @Produce json
// @Param group path string true "Group ID"
// @Success 200 {array} entity.Channel "Channels retrieved successfully"
// @Failure 400 {object} map[string]string "Group ID is required"
// @Failure 404 {object} map[string]string "Channels not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /channel/group/{group} [get]
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

	channel, err := usecase.ListByGroup(groupID)
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

// FindByPlatform godoc
// @Summary Get channels by platform
// @Description Retrieves a list of channels associated with a specific platform ID.
// @Tags channel
// @Produce json
// @Param platform path string true "Platform ID"
// @Success 200 {array} entity.Channel "Channels retrieved successfully"
// @Failure 400 {object} map[string]string "Platform ID is required"
// @Failure 404 {object} map[string]string "Channels not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /channel/platform/{platform} [get]
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

	channel, err := usecase.ListByPlatform(platformID)
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

// DeleteById godoc
// @Summary Delete channel by ID
// @Description Deletes a channel by its unique identifier.
// @Tags channel
// @Produce json
// @Param id path string true "Channel ID"
// @Success 200 {string} string "Channel deleted successfully"
// @Failure 400 {object} map[string]string "Channel ID is required"
// @Failure 404 {object} map[string]string "Channel not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /channel/{id} [delete]
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
