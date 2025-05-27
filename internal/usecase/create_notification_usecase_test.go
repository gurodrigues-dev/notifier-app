package usecase

import (
	"errors"
	"testing"

	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/value"
	"github.com/gurodrigues-dev/notifier-app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateNotificationUsecase_CreateNotification(t *testing.T) {
	type args struct {
		input value.NotificationInput
	}
	type fields struct {
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *CreateNotificationUsecase
		wantErr bool
	}{
		{
			name: "there is to return success",
			args: args{
				input: value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Title:    "Payment Success",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				},
			},
			setup: func(t *testing.T) *CreateNotificationUsecase {
				input := value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				}
				notificationRepository := mocks.NewNotificationRepository(t)
				channelRepository := mocks.NewChannelRepository(t)
				cacher := mocks.NewCacher(t)
				kafka := mocks.NewQueue(t)
				logger := mocks.NewLogger(t)

				logger.On("Infof", mock.Anything, mock.Anything).Return()
				channelRepository.On("GetByIDs", []string{"2"}).Return([]entity.Channel{}, nil)
				channelRepository.On("GetByGroups", []string{"marketing"}).Return([]entity.Channel{}, nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				cacher.On("Get", input.UUID).Return("", nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				cacher.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				kafka.On("Produce", mock.Anything, mock.Anything).Return(nil)

				return NewCreateNotificationUsecase(
					notificationRepository,
					channelRepository,
					cacher,
					kafka,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is to return fail to get channels by ids",
			args: args{
				input: value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Title:    "Payment Success",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				},
			},
			setup: func(t *testing.T) *CreateNotificationUsecase {
				notificationRepository := mocks.NewNotificationRepository(t)
				channelRepository := mocks.NewChannelRepository(t)
				cacher := mocks.NewCacher(t)
				kafka := mocks.NewQueue(t)
				logger := mocks.NewLogger(t)

				logger.On("Infof", mock.Anything, mock.Anything).Return()
				channelRepository.On("GetByIDs", []string{"2"}).Return([]entity.Channel{}, errors.New("get by id db error"))
				channelRepository.On("GetByGroups", []string{"marketing"}).Return([]entity.Channel{}, nil)

				notificationRepository.On("CreateNotification", mock.Anything).Maybe()

				return NewCreateNotificationUsecase(
					notificationRepository,
					channelRepository,
					cacher,
					kafka,
					logger,
				)
			},
			wantErr: true,
		},
		{
			name: "there is to return has cache",
			args: args{
				input: value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Title:    "Payment Success",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				},
			},
			setup: func(t *testing.T) *CreateNotificationUsecase {
				input := value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				}
				notificationRepository := mocks.NewNotificationRepository(t)
				channelRepository := mocks.NewChannelRepository(t)
				cacher := mocks.NewCacher(t)
				kafka := mocks.NewQueue(t)
				logger := mocks.NewLogger(t)

				logger.On("Infof", mock.Anything, mock.Anything).Return()
				channelRepository.On("GetByIDs", []string{"2"}).Return([]entity.Channel{}, nil)
				channelRepository.On("GetByGroups", []string{"marketing"}).Return([]entity.Channel{}, nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				cacher.On("Get", input.UUID).Return("cache", nil)

				notificationRepository.On("CreateNotification", mock.Anything).Maybe().Return(nil)

				return NewCreateNotificationUsecase(
					notificationRepository,
					channelRepository,
					cacher,
					kafka,
					logger,
				)
			},
			wantErr: true,
		},
		{
			name: "there is to return error cache",
			args: args{
				input: value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Title:    "Payment Success",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				},
			},
			setup: func(t *testing.T) *CreateNotificationUsecase {
				input := value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				}
				notificationRepository := mocks.NewNotificationRepository(t)
				channelRepository := mocks.NewChannelRepository(t)
				cacher := mocks.NewCacher(t)
				kafka := mocks.NewQueue(t)
				logger := mocks.NewLogger(t)

				logger.On("Infof", mock.Anything, mock.Anything).Return()
				channelRepository.On("GetByIDs", []string{"2"}).Return([]entity.Channel{}, nil)
				channelRepository.On("GetByGroups", []string{"marketing"}).Return([]entity.Channel{}, nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				cacher.On("Get", input.UUID).Return("", errors.New("get cache error"))

				notificationRepository.On("CreateNotification", mock.Anything).Maybe().Return(nil)

				return NewCreateNotificationUsecase(
					notificationRepository,
					channelRepository,
					cacher,
					kafka,
					logger,
				)
			},
			wantErr: true,
		},
		{
			name: "there is to return error to set cache",
			args: args{
				input: value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Title:    "Payment Success",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				},
			},
			setup: func(t *testing.T) *CreateNotificationUsecase {
				input := value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				}
				notificationRepository := mocks.NewNotificationRepository(t)
				channelRepository := mocks.NewChannelRepository(t)
				cacher := mocks.NewCacher(t)
				kafka := mocks.NewQueue(t)
				logger := mocks.NewLogger(t)

				logger.On("Infof", mock.Anything, mock.Anything).Return()
				channelRepository.On("GetByIDs", []string{"2"}).Return([]entity.Channel{}, nil)
				channelRepository.On("GetByGroups", []string{"marketing"}).Return([]entity.Channel{}, nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				cacher.On("Get", input.UUID).Return("", nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				cacher.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("set cache error"))

				notificationRepository.On("CreateNotification", mock.Anything).Maybe().Return(nil)

				return NewCreateNotificationUsecase(
					notificationRepository,
					channelRepository,
					cacher,
					kafka,
					logger,
				)
			},
			wantErr: true,
		},
		{
			name: "there is to return error to produce message to kafka",
			args: args{
				input: value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Title:    "Payment Success",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				},
			},
			setup: func(t *testing.T) *CreateNotificationUsecase {
				input := value.NotificationInput{
					Channels: []string{"marketing", "2"},
					UUID:     "2bbcdd20-1ea6-42be-8484-02f3007e3463",
					Event: value.Event{
						Name:      "payment_success",
						Timestamp: 1748355999,
						Requester: "requester",
						Receiver:  "receiver",
						Currency:  "BRL",
						Category:  "pix",
						CostCents: 9000,
					},
				}
				notificationRepository := mocks.NewNotificationRepository(t)
				channelRepository := mocks.NewChannelRepository(t)
				cacher := mocks.NewCacher(t)
				kafka := mocks.NewQueue(t)
				logger := mocks.NewLogger(t)

				logger.On("Infof", mock.Anything, mock.Anything).Return()
				channelRepository.On("GetByIDs", []string{"2"}).Return([]entity.Channel{}, nil)
				channelRepository.On("GetByGroups", []string{"marketing"}).Return([]entity.Channel{}, nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				cacher.On("Get", input.UUID).Return("", nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				cacher.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				kafka.On("Produce", mock.Anything, mock.Anything).Return(errors.New("kafka produce error"))

				notificationRepository.On("CreateNotification", mock.Anything).Maybe().Return(nil)

				return NewCreateNotificationUsecase(
					notificationRepository,
					channelRepository,
					cacher,
					kafka,
					logger,
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := tt.setup(t)
			err := usecase.CreateNotification(tt.args.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
