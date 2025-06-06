package usecase

import (
	"errors"
	"testing"

	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTokenUsecase_CreateToken(t *testing.T) {
	type args struct {
		token *entity.Token
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *CreateTokenUsecase
		wantErr bool
	}{
		{
			name: "there is to return success",
			args: args{
				token: &entity.Token{},
			},
			setup: func(t *testing.T) *CreateTokenUsecase {
				repository := mocks.NewAuthRepository(t)
				logger := mocks.NewLogger(t)
				repository.On("CreateToken", mock.Anything).Return(nil)
				logger.On("Infof", mock.Anything, mock.Anything).Return()
				return NewCreateTokenUsecase(
					repository,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is to return db error",
			args: args{
				token: &entity.Token{},
			},
			setup: func(t *testing.T) *CreateTokenUsecase {
				repository := mocks.NewAuthRepository(t)
				logger := mocks.NewLogger(t)
				repository.On("CreateToken", mock.Anything).Return(errors.New("db error"))
				return NewCreateTokenUsecase(
					repository,
					logger,
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := tt.setup(t)
			_, err := usecase.CreateToken(tt.args.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
