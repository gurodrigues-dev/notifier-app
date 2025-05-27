package usecase

import (
	"fmt"

	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
	"github.com/gurodrigues-dev/notifier-app/internal/value"
	"github.com/gurodrigues-dev/notifier-app/pkg/slicecommon"
)

type CreateChannelUsecase struct {
	channelRepository repository.ChannelRepository
	ses               contracts.SESIface
	logger            contracts.Logger
}

func NewCreateChannelUsecase(
	channelRepository repository.ChannelRepository,
	ses contracts.SESIface,
	logger contracts.Logger,
) *CreateChannelUsecase {
	return &CreateChannelUsecase{
		channelRepository: channelRepository,
		ses:               ses,
		logger:            logger,
	}
}

func (ccu *CreateChannelUsecase) CreateChannel(channel *entity.Channel) (*entity.Channel, error) {
	if slicecommon.Contains(value.Platforms, channel.Platform) == false {
		return nil, fmt.Errorf("invalid plataform: %s", channel.Platform)
	}
	if channel.Platform == value.EmailPlatform {
		err := ccu.ses.VerifyEmail(channel.TargetID)
		if err != nil {
			return nil, err
		}
	}
	return ccu.channelRepository.CreateChannel(channel)
}
