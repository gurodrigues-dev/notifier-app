package usecase

import (
	"errors"
	"testing"

	"github.com/gurodrigues-dev/notifier-app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTokenUsecase_DeleteToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *DeleteTokenUsecase
		wantErr bool
	}{
		{
			name: "there is to return success",
			args: args{
				token: "valid-token",
			},
			setup: func(t *testing.T) *DeleteTokenUsecase {
				mock := mocks.NewAuthRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("DeleteToken", "valid-token").Return(nil)
				return NewDeleteTokenUsecase(
					mock,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is to return db error",
			args: args{
				token: "valid-token",
			},
			setup: func(t *testing.T) *DeleteTokenUsecase {
				mock := mocks.NewAuthRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("DeleteToken", "valid-token").Return(errors.New("db error"))
				return NewDeleteTokenUsecase(
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
			err := usecase.DeleteToken(tt.args.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
