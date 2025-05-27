package persistence

import (
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type NotificationRepositoryImpl struct {
	Postgres contracts.PostgresIface
}

func (nr NotificationRepositoryImpl) CreateNotification(notification *entity.NotificationError) error {
	return nr.Postgres.Client().Create(notification).Error
}

func (nr NotificationRepositoryImpl) GetNotificationByID(id string) (*entity.NotificationError, error) {
	var notification entity.NotificationError
	err := nr.Postgres.Client().Where("id = ?", id).First(&notification).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}
