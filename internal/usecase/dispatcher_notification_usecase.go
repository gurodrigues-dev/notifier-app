package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
	"github.com/gurodrigues-dev/notifier-app/internal/value"
	"github.com/gurodrigues-dev/notifier-app/pkg/stringcommon"
)

type DispatcherUsecase struct {
	notificationRepository repository.NotificationRepository
	ses                    contracts.SESIface
	logger                 contracts.Logger
}

func NewDispatcherUsecase(
	notificationRepository repository.NotificationRepository,
	ses contracts.SESIface,
	logger contracts.Logger,
) *DispatcherUsecase {
	return &DispatcherUsecase{
		notificationRepository: notificationRepository,
		ses:                    ses,
		logger:                 logger,
	}
}

func (du *DispatcherUsecase) Execute(message string) (err error) {
	var successfullyNotifications []string
	var errorNotifications []string
	var notification *entity.Notification
	err = json.Unmarshal([]byte(message), &notification)
	if err != nil {
		return err
	}

	for _, channel := range notification.Channels {
		if channel.Platform == value.EmailPlatform {
			err := du.ses.SendEmail(
				&entity.Email{
					Recipient: channel.TargetID,
					Subject:   notification.Title,
					Body:      notification.Message,
				},
			)
			if err != nil {
				errorNotifications = append(errorNotifications, err.Error())
				continue
			}
			successfullyNotifications = append(successfullyNotifications, notification.UUID)
		}

		if channel.Platform == value.DiscordPlatform || channel.Platform == value.SlackPlatform {
			err := du.httpCall(channel.TargetID, channel.Platform, fmt.Sprintf("%s: %s", notification.Title, notification.Message))
			if err != nil {
				errorNotifications = append(errorNotifications, err.Error())
				continue
			}
			successfullyNotifications = append(successfullyNotifications, notification.UUID)
		}
	}

	if len(successfullyNotifications) > 0 {
		return nil
	}

	if notification.Retries < value.MaxRetries {
		notification.Retries++
		serializedMessage, serialError := stringcommon.SerializeToJSON(notification)
		if serialError != nil {
			return fmt.Errorf("error serializing notifcation: %w", err)
		}

		dbError := du.notificationRepository.CreateNotification(&entity.NotificationError{
			UUID:  notification.UUID,
			Body:  serializedMessage,
			Error: strings.Join(errorNotifications, ", "),
		})
		if dbError != nil {
			du.logger.Errorf("error creating notification")
			return dbError
		}
	}

	return fmt.Errorf("")
}

func (du *DispatcherUsecase) httpCall(webhook, platform, message string) error {
	payload := make(map[string]string)

	if platform == value.DiscordPlatform {
		payload["content"] = message
	}

	if platform == value.SlackPlatform {
		payload["text"] = message
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("slack webhook returned status code %d", resp.StatusCode)
	}

	return nil
}
