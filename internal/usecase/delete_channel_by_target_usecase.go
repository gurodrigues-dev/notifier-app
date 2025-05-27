package usecase

import (
	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type DeleteChannelByIDUsecase struct {
	channelRepository repository.ChannelRepository
	logger            contracts.Logger
}

func NewDeleteChannelByIDUsecase(
	channelRepository repository.ChannelRepository,
	logger contracts.Logger,
) *DeleteChannelByIDUsecase {
	return &DeleteChannelByIDUsecase{
		channelRepository: channelRepository,
		logger:            logger,
	}
}

func (dcu *DeleteChannelByIDUsecase) DeleteByID(targetID string) error {
	return dcu.channelRepository.DeleteByID(targetID)
}
