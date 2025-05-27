package usecase

import (
	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type CreateTokenUsecase struct {
	authRepository repository.AuthRepository
	logger         contracts.Logger
}

func NewCreateTokenUsecase(
	authRepository repository.AuthRepository,
	logger contracts.Logger,
) *CreateTokenUsecase {
	return &CreateTokenUsecase{
		authRepository: authRepository,
		logger:         logger,
	}
}

func (ctu *CreateTokenUsecase) CreateToken(token *entity.Token) (string, error) {
	token.Token = token.CreateToken()

	err := ctu.authRepository.CreateToken(token)
	if err != nil {
		return "", err
	}

	ctu.logger.Infof("Token created")
	return token.Token, nil
}
