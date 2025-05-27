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

// CreateNotification godoc
// @Summary Create a new notification
// @Description Creates a new notification based on the provided notification data.
// @Tags notification
// @Accept json
// @Produce json
// @Param notification body value.NotificationInput true "Notification request body"
// @Success 200 {string} string "Notification sent successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 422 {object} map[string]string "Unprocessable entity"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notification [post]
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

// GetNotification godoc
// @Summary Get notification by ID
// @Description Retrieves a notification by its unique identifier.
// @Tags notification
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} value.NotificationInput "Notification retrieved successfully"
// @Failure 404 {object} map[string]string "Notification not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notification/{id} [get]
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
