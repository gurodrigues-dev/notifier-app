package usecase

import (
	"encoding/json"

	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
	"github.com/gurodrigues-dev/notifier-app/internal/value"
)

type GetNotificationUsecase struct {
	notificationRepository repository.NotificationRepository
	logger                 contracts.Logger
}

func NewGetNotificationUsecase(
	notificationRepository repository.NotificationRepository,
	logger contracts.Logger,
) *GetNotificationUsecase {
	return &GetNotificationUsecase{
		notificationRepository: notificationRepository,
		logger:                 logger,
	}
}

func (gnu *GetNotificationUsecase) GetNotification(id string) (*value.NotificationOutput, error) {
	var body entity.Notification
	notification, err := gnu.notificationRepository.GetNotificationByID(id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(notification.Body, &body)
	if err != nil {
		return nil, err
	}

	return &value.NotificationOutput{
		ID:    notification.ID,
		UUID:  notification.UUID,
		Body:  body,
		Error: notification.Error,
	}, nil
}
