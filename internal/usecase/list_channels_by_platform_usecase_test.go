package usecase

import (
	"errors"
	"testing"

	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestListChannelsByPlatformUsecase_ListByPlatform(t *testing.T) {
	type args struct {
		platform string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *ListChannelsByPlatformUsecase
		wantErr bool
	}{
		{
			name: "there is to return success",
			args: args{
				platform: "slack",
			},
			setup: func(t *testing.T) *ListChannelsByPlatformUsecase {
				mock := mocks.NewChannelRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("GetByPlatform", "slack").Return([]entity.Channel{}, nil)
				return NewListChannelsByPlatformUsecase(
					mock,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is to return db error",
			args: args{
				platform: "slack",
			},
			setup: func(t *testing.T) *ListChannelsByPlatformUsecase {
				mock := mocks.NewChannelRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("GetByPlatform", "slack").Return([]entity.Channel{}, errors.New("db error"))
				return NewListChannelsByPlatformUsecase(
					mock,
					logger,
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := tt.setup(t)
			_, err := usecase.ListByPlatform(tt.args.platform)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
