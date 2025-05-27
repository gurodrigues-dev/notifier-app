package usecase

import (
	"errors"
	"testing"

	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetChannelByIDUsecase_GetByID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *GetChannelByIDUsecase
		wantErr bool
	}{
		{
			name: "there is to return success",
			args: args{
				id: "1",
			},
			setup: func(t *testing.T) *GetChannelByIDUsecase {
				mock := mocks.NewChannelRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("GetByID", "1").Return(&entity.Channel{}, nil)
				return NewGetChannelByIDUsecase(
					mock,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is to return db error",
			args: args{
				id: "1",
			},
			setup: func(t *testing.T) *GetChannelByIDUsecase {
				mock := mocks.NewChannelRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("GetByID", "1").Return(&entity.Channel{}, errors.New("db error"))
				return NewGetChannelByIDUsecase(
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
			_, err := usecase.GetByID(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
