package usecase

import (
	"errors"
	"testing"

	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetTokenUsecase_GetToken(t *testing.T) {
	type args struct {
		user string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *GetTokenByUserUsecase
		wantErr bool
	}{
		{
			name: "there is to return success",
			args: args{
				user: "test_user",
			},
			setup: func(t *testing.T) *GetTokenByUserUsecase {
				mock := mocks.NewAuthRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("GetTokenByUser", "test_user").Return(&entity.Token{
					Token: "test_token123abc",
				}, nil)
				return NewGetTokenByUserUsecase(
					mock,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is to return db error",
			args: args{
				user: "test_user",
			},
			setup: func(t *testing.T) *GetTokenByUserUsecase {
				mock := mocks.NewAuthRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("GetTokenByUser", "test_user").Return(&entity.Token{}, errors.New("db error"))
				return NewGetTokenByUserUsecase(
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
			_, err := usecase.GetTokenByUser(tt.args.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
