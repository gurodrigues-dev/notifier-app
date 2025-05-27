package repository

import (
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
)

type NotificationRepository interface {
	CreateNotification(notification *entity.NotificationError) error
	GetNotificationByID(id string) (*entity.NotificationError, error)
}
