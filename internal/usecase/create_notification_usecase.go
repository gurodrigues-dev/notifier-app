package usecase

import (
	"fmt"
	"sync"
	"time"

	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
	"github.com/gurodrigues-dev/notifier-app/internal/value"
	"github.com/gurodrigues-dev/notifier-app/pkg/slicecommon"
	"github.com/gurodrigues-dev/notifier-app/pkg/stringcommon"
	"github.com/redis/go-redis/v9"
)

type CreateNotificationUsecase struct {
	notificationRepository repository.NotificationRepository
	channelRepository      repository.ChannelRepository
	cacher                 contracts.Cacher
	queue                  contracts.Queue
	logger                 contracts.Logger
}

func NewCreateNotificationUsecase(
	notificationRepository repository.NotificationRepository,
	channelRepository repository.ChannelRepository,
	cacher contracts.Cacher,
	queue contracts.Queue,
	logger contracts.Logger,
) *CreateNotificationUsecase {
	return &CreateNotificationUsecase{
		notificationRepository: notificationRepository,
		channelRepository:      channelRepository,
		cacher:                 cacher,
		queue:                  queue,
		logger:                 logger,
	}
}

var (
	timeDuration = 2 * time.Hour
)

func (cnu *CreateNotificationUsecase) CreateNotification(input value.NotificationInput) (err error) {
	cnu.logger.Infof("getting channels")
	channels, err := cnu.GetChannels(input)
	if err != nil {
		return err
	}

	cnu.logger.Infof("build notification")
	notification := cnu.buildNotificationMessage(input)
	notification.Channels = channels

	cnu.logger.Infof("serializing notification")
	serializedMessage, err := stringcommon.SerializeToJSON(notification)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			dbError := cnu.notificationRepository.CreateNotification(&entity.NotificationError{
				UUID:  notification.UUID,
				Body:  serializedMessage,
				Error: err.Error(),
			})
			if dbError != nil {
				cnu.logger.Errorf("error creating notification")
				return
			}
		}
	}()

	cnu.logger.Infof("validating cache")
	hasCache, err := cnu.isCacheValid(notification.UUID)
	if err != nil && err != redis.Nil {
		return err
	}

	if hasCache {
		return fmt.Errorf("notification already exists in cache")
	}

	cnu.logger.Infof("creating cache")
	err = cnu.setCache(notification.UUID, *notification)
	if err != nil {
		return err
	}

	cnu.logger.Infof("sending notification")
	err = cnu.queue.Produce(value.GetTopic(), string(serializedMessage))
	if err != nil {
		return err
	}

	return nil
}

func (cnu *CreateNotificationUsecase) GetChannels(input value.NotificationInput) (map[int]entity.Channel, error) {
	target, group := slicecommon.Partition(input.Channels)

	var (
		wg       sync.WaitGroup
		channels sync.Map
		errCh    = make(chan error, 2)
	)

	if len(target) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			targets, err := cnu.channelRepository.GetByIDs(target)
			if err != nil {
				errCh <- err
				return
			}
			for _, t := range targets {
				channels.Store(t.ID, t)
			}
		}()
	}

	if len(group) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			groups, err := cnu.channelRepository.GetByGroups(group)
			if err != nil {
				errCh <- err
				return
			}
			for _, g := range groups {
				channels.Store(g.ID, g)
			}
		}()
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	result := make(map[int]entity.Channel)
	channels.Range(func(key, value any) bool {
		id := key.(int)
		channel := value.(entity.Channel)
		result[id] = channel
		return true
	})

	return result, nil
}

func (cnu *CreateNotificationUsecase) buildNotificationMessage(input value.NotificationInput) *entity.Notification {
	timestamp := time.Unix(input.Event.Timestamp, 0)
	formattedTime := timestamp.Format("02/01/06 15:04")

	messageText := fmt.Sprintf(
		"There was a new transaction at %s, between %s and %s by %s, with the value of %s %v, status: %s",
		formattedTime,
		input.Event.Requester,
		input.Event.Receiver,
		input.Event.Category,
		input.Event.Currency,
		input.Event.CostCents/100,
		input.Event.Name,
	)

	notification := &entity.Notification{
		UUID:    input.UUID,
		Title:   input.Title,
		Message: messageText,
		Event: entity.Event{
			Name:      input.Event.Name,
			Currency:  input.Event.Currency,
			Requester: input.Event.Requester,
			Receiver:  input.Event.Receiver,
			Category:  input.Event.Category,
			Timestamp: input.Event.Timestamp,
			CostCents: input.Event.CostCents,
		},
	}

	return notification
}

func (cnu *CreateNotificationUsecase) isCacheValid(uuid string) (bool, error) {
	value, err := cnu.cacher.Get(uuid)
	if err != nil {
		return false, err
	}
	return value != "", nil
}

func (cnu *CreateNotificationUsecase) setCache(uuid string, value entity.Notification) error {
	err := cnu.cacher.Set(uuid, value, timeDuration)
	if err != nil {
		return err
	}

	return nil
}
