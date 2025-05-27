package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gurodrigues-dev/notifier-app/internal/controller"
	"github.com/gurodrigues-dev/notifier-app/internal/domain/middleware"
	"github.com/gurodrigues-dev/notifier-app/internal/infra"
)

type Controllers struct {
	Auth         *controller.AuthController
	Notification *controller.NotificationController
	Channel      *controller.ChannelController
}

func NewControllers() *Controllers {
	return &Controllers{
		Notification: controller.NewNotificationController(infra.App.Logger),
		Auth:         controller.NewAuthController(infra.App.Logger),
		Channel:      controller.NewChannelController(infra.App.Logger),
	}
}

func (routes *Controllers) Routes(group *gin.RouterGroup, middleware *middleware.Middleware) {
	group.POST("/notification", middleware.TokenMiddleware(), routes.Notification.CreateNotification)
	group.GET("/notification/:id", middleware.AdminMiddleware(), routes.Notification.GetNotification)

	group.POST("/token", middleware.AdminMiddleware(), routes.Auth.CreateToken)
	group.GET("/token/:user", middleware.AdminMiddleware(), routes.Auth.GetToken)
	group.DELETE("/token/:user", middleware.AdminMiddleware(), routes.Auth.DeleteToken)

	group.POST("/channel", middleware.TokenMiddleware(), routes.Channel.CreateChannel)
	group.GET("/channel/:id", middleware.TokenMiddleware(), routes.Channel.FindById)
	group.DELETE("/channel/:id", middleware.TokenMiddleware(), routes.Channel.DeleteById)
	group.GET("/group/:group", middleware.TokenMiddleware(), routes.Channel.FindByGroup)
	group.GET("/platform/:platform", middleware.TokenMiddleware(), routes.Channel.FindByPlatform)
}
