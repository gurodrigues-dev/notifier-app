package usecase

import (
	"errors"
	"testing"

	"github.com/gurodrigues-dev/notifier-app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteChannelUsecase_DeleteChannel(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *DeleteChannelByIDUsecase
		wantErr bool
	}{
		{
			name: "there is to return success",
			args: args{
				id: "1",
			},
			setup: func(t *testing.T) *DeleteChannelByIDUsecase {
				mock := mocks.NewChannelRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("DeleteByID", "1").Return(nil)
				return NewDeleteChannelByIDUsecase(
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
			setup: func(t *testing.T) *DeleteChannelByIDUsecase {
				mock := mocks.NewChannelRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("DeleteByID", "1").Return(errors.New("db error"))
				return NewDeleteChannelByIDUsecase(
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
			err := usecase.DeleteByID(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
