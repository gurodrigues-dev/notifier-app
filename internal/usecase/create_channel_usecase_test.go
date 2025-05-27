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

func TestCreateChannelUsecase_CreateChannel(t *testing.T) {
	type args struct {
		channel *entity.Channel
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *CreateChannelUsecase
		wantErr bool
	}{
		{
			name: "there is to return success",
			args: args{
				channel: &entity.Channel{
					Platform: value.SlackPlatform,
				},
			},
			setup: func(t *testing.T) *CreateChannelUsecase {
				repository := mocks.NewChannelRepository(t)
				ses := mocks.NewSESIface(t)
				logger := mocks.NewLogger(t)
				repository.On("CreateChannel", &entity.Channel{
					Platform: value.SlackPlatform,
				}).Return(&entity.Channel{}, nil)
				return NewCreateChannelUsecase(
					repository,
					ses,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is to return success using email platform",
			args: args{
				channel: &entity.Channel{
					Platform: value.EmailPlatform,
				},
			},
			setup: func(t *testing.T) *CreateChannelUsecase {
				repository := mocks.NewChannelRepository(t)
				ses := mocks.NewSESIface(t)
				logger := mocks.NewLogger(t)
				repository.On("CreateChannel", &entity.Channel{
					Platform: value.EmailPlatform,
				}).Return(&entity.Channel{}, nil)
				ses.On("VerifyEmail", mock.Anything).Return(nil)
				return NewCreateChannelUsecase(
					repository,
					ses,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is to return platform error",
			args: args{
				channel: &entity.Channel{
					Platform: "nonexistingplatform",
				},
			},
			setup: func(t *testing.T) *CreateChannelUsecase {
				repository := mocks.NewChannelRepository(t)
				ses := mocks.NewSESIface(t)
				logger := mocks.NewLogger(t)
				return NewCreateChannelUsecase(
					repository,
					ses,
					logger,
				)
			},
			wantErr: true,
		},
		{
			name: "there is to return db error",
			args: args{
				channel: &entity.Channel{Platform: value.SlackPlatform},
			},
			setup: func(t *testing.T) *CreateChannelUsecase {
				repository := mocks.NewChannelRepository(t)
				ses := mocks.NewSESIface(t)
				logger := mocks.NewLogger(t)
				repository.On("CreateChannel", &entity.Channel{Platform: value.SlackPlatform}).Return(&entity.Channel{}, errors.New("db error"))
				return NewCreateChannelUsecase(
					repository,
					ses,
					logger,
				)
			},
			wantErr: true,
		},
		{
			name: "there is to return error using email platform",
			args: args{
				channel: &entity.Channel{
					Platform: value.EmailPlatform,
				},
			},
			setup: func(t *testing.T) *CreateChannelUsecase {
				repository := mocks.NewChannelRepository(t)
				ses := mocks.NewSESIface(t)
				logger := mocks.NewLogger(t)
				ses.On("VerifyEmail", mock.Anything).Return(errors.New("email error"))
				return NewCreateChannelUsecase(
					repository,
					ses,
					logger,
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := tt.setup(t)
			_, err := usecase.CreateChannel(tt.args.channel)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
