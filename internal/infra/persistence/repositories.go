package persistence

import "github.com/gurodrigues-dev/notifier-app/internal/domain/repository"

type PostgresRepositories struct {
	AuthRepository         repository.AuthRepository
	NotificationRepository repository.NotificationRepository
	ChannelRepository      repository.ChannelRepository
}
