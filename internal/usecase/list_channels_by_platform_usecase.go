package usecase

import (
	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type ListChannelsByPlatformUsecase struct {
	channelRepository repository.ChannelRepository
	logger            contracts.Logger
}

func NewListChannelsByPlatformUsecase(
	channelRepository repository.ChannelRepository,
	logger contracts.Logger,
) *ListChannelsByPlatformUsecase {
	return &ListChannelsByPlatformUsecase{
		channelRepository: channelRepository,
		logger:            logger,
	}
}

func (lcu *ListChannelsByPlatformUsecase) ListByPlatform(platform string) ([]entity.Channel, error) {
	return lcu.channelRepository.GetByPlatform(platform)
}
