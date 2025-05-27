package infra

import (
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/persistence"
)

type Application struct {
	Repositories persistence.PostgresRepositories
	Postgres     contracts.PostgresIface
	Cache        contracts.Cacher
	Email        contracts.SESIface
	Logger       contracts.Logger
	Queue        contracts.Queue
	Metrics      contracts.Metrics
}

var App Application
