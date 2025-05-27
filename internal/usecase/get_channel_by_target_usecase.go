package usecase

import (
	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type GetChannelByIDUsecase struct {
	channelRepository repository.ChannelRepository
	logger            contracts.Logger
}

func NewGetChannelByIDUsecase(
	channelRepository repository.ChannelRepository,
	logger contracts.Logger,
) *GetChannelByIDUsecase {
	return &GetChannelByIDUsecase{
		channelRepository: channelRepository,
		logger:            logger,
	}
}

func (gcu *GetChannelByIDUsecase) GetByID(targetID string) (*entity.Channel, error) {
	return gcu.channelRepository.GetByID(targetID)
}
