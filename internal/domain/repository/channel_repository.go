package repository

import "github.com/gurodrigues-dev/notifier-app/internal/entity"

type ChannelRepository interface {
	CreateChannel(channel *entity.Channel) (*entity.Channel, error)
	GetByID(id string) (*entity.Channel, error)
	GetByIDs(ids []string) ([]entity.Channel, error)
	GetByGroup(group string) ([]entity.Channel, error)
	GetByGroups(groups []string) ([]entity.Channel, error)
	GetByPlatform(platform string) ([]entity.Channel, error)
	DeleteByID(id string) error
}
