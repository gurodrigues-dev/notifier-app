package usecase

import (
	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type GetTokenByUserUsecase struct {
	authRepository repository.AuthRepository
	logger         contracts.Logger
}

func NewGetTokenByUserUsecase(
	authRepository repository.AuthRepository,
	logger contracts.Logger,
) *GetTokenByUserUsecase {
	return &GetTokenByUserUsecase{
		authRepository: authRepository,
		logger:         logger,
	}
}

func (gtu *GetTokenByUserUsecase) GetTokenByUser(user string) (*entity.Token, error) {
	return gtu.authRepository.GetTokenByUser(user)
}
