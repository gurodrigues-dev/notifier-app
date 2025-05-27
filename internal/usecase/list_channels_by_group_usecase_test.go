package usecase

import (
	"errors"
	"testing"

	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestListChannelsByGroupUsecase_ListByGroupID(t *testing.T) {
	type args struct {
		group string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *ListChannelsByGroupUsecase
		wantErr bool
	}{
		{
			name: "there is to return success",
			args: args{
				group: "marketing",
			},
			setup: func(t *testing.T) *ListChannelsByGroupUsecase {
				mock := mocks.NewChannelRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("GetByGroup", "marketing").Return([]entity.Channel{}, nil)
				return NewListChannelsByGroupUsecase(
					mock,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is to return db error",
			args: args{
				group: "marketing",
			},
			setup: func(t *testing.T) *ListChannelsByGroupUsecase {
				mock := mocks.NewChannelRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("GetByGroup", "marketing").Return([]entity.Channel{}, errors.New("db error"))
				return NewListChannelsByGroupUsecase(
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
			_, err := usecase.ListByGroup(tt.args.group)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
