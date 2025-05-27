package usecase

import (
	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type DeleteTokenUsecase struct {
	authRepository repository.AuthRepository
	logger         contracts.Logger
}

func NewDeleteTokenUsecase(
	authRepository repository.AuthRepository,
	logger contracts.Logger,
) *DeleteTokenUsecase {
	return &DeleteTokenUsecase{
		authRepository: authRepository,
		logger:         logger,
	}
}

func (dtu *DeleteTokenUsecase) DeleteToken(token string) error {
	return dtu.authRepository.DeleteToken(token)
}
