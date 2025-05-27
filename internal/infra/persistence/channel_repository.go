package persistence

import (
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type ChannelRepositoryImpl struct {
	Postgres contracts.PostgresIface
}

func (cr ChannelRepositoryImpl) CreateChannel(channel *entity.Channel) (*entity.Channel, error) {
	if err := cr.Postgres.Client().Create(channel).Error; err != nil {
		return nil, err
	}
	return channel, nil
}

func (cr ChannelRepositoryImpl) GetByID(id string) (*entity.Channel, error) {
	var channel entity.Channel
	err := cr.Postgres.Client().Where("id = ?", id).First(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (cr ChannelRepositoryImpl) GetByIDs(ids []string) ([]entity.Channel, error) {
	var channels []entity.Channel
	err := cr.Postgres.Client().Where("id IN (?)", ids).Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (cr ChannelRepositoryImpl) GetByGroup(group string) ([]entity.Channel, error) {
	var channels []entity.Channel
	err := cr.Postgres.Client().Where(`"group" = ?`, group).Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (cr ChannelRepositoryImpl) GetByGroups(groups []string) ([]entity.Channel, error) {
	var channels []entity.Channel
	err := cr.Postgres.Client().Where(`"group" IN (?)`, groups).Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (cr ChannelRepositoryImpl) GetByPlatform(platform string) ([]entity.Channel, error) {
	var channels []entity.Channel
	err := cr.Postgres.Client().Where("platform = ?", platform).Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (cr ChannelRepositoryImpl) DeleteByID(id string) error {
	return cr.Postgres.Client().Where("id = ?", id).Delete(&entity.Channel{}).Error
}
