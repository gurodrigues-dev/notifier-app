package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gurodrigues-dev/notifier-app/internal/infra"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
	"github.com/gurodrigues-dev/notifier-app/internal/usecase"
	"github.com/gurodrigues-dev/notifier-app/internal/value"
	"github.com/jinzhu/gorm"
)

type NotificationController struct {
	logger contracts.Logger
}

func NewNotificationController(
	logger contracts.Logger,
) *NotificationController {
	return &NotificationController{
		logger: logger,
	}
}

func (nc *NotificationController) CreateNotification(httpContext *gin.Context) {
	var requestParams value.NotificationInput
	if err := httpContext.BindJSON(&requestParams); err != nil {
		httpContext.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestParams); err != nil {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecase := usecase.NewCreateNotificationUsecase(
		infra.App.Repositories.NotificationRepository,
		infra.App.Repositories.ChannelRepository,
		infra.App.Cache,
		infra.App.Queue,
		nc.logger,
	)

	err := usecase.CreateNotification(requestParams)
	if err != nil {
		httpContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	httpContext.JSON(http.StatusOK, "notification sent successfully")
}

func (nc *NotificationController) GetNotification(httpContext *gin.Context) {
	id := httpContext.Param("id")

	usecase := usecase.NewGetNotificationUsecase(
		infra.App.Repositories.NotificationRepository,
		nc.logger,
	)

	notifyError, err := usecase.GetNotification(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httpContext.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpContext.JSON(http.StatusOK, notifyError)
	return
}
