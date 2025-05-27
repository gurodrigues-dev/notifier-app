package usecase

import (
	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type ListChannelsByGroupUsecase struct {
	channelRepository repository.ChannelRepository
	logger            contracts.Logger
}

func NewListChannelsByGroupUsecase(
	channelRepository repository.ChannelRepository,
	logger contracts.Logger,
) *ListChannelsByGroupUsecase {
	return &ListChannelsByGroupUsecase{
		channelRepository: channelRepository,
		logger:            logger,
	}
}

func (lcu *ListChannelsByGroupUsecase) ListByGroupID(group string) ([]entity.Channel, error) {
	return lcu.channelRepository.GetByGroup(group)
}
