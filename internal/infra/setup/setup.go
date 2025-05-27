package setup

import (
	"log"
	"os"

	"github.com/gurodrigues-dev/notifier-app/config"
	"github.com/gurodrigues-dev/notifier-app/internal/infra"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/cache"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/database"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/email"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/logger"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/persistence"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/queue"
	"github.com/gurodrigues-dev/notifier-app/internal/metrics"
	"github.com/spf13/viper"
)

const (
	ServiceName = "notifier"
)

type Setup struct {
	app          *infra.Application
	repositories *persistence.PostgresRepositories
}

func NewSetup() Setup {
	err := config.LoadServerEnvironmentVars(ServiceName, os.Getenv(config.ServerEnvironment))
	if err != nil {
		log.Fatal(err)
	}

	return Setup{
		app:          &infra.Application{},
		repositories: &persistence.PostgresRepositories{},
	}
}

func (s Setup) Postgres() {
	s.app.Postgres, _ = database.NewPGGORMImpl()
}

func (s Setup) Repositories() {
	s.app.Repositories.NotificationRepository = persistence.NotificationRepositoryImpl{Postgres: s.app.Postgres}
	s.app.Repositories.AuthRepository = persistence.AuthRepositoryImpl{Postgres: s.app.Postgres}
	s.app.Repositories.ChannelRepository = persistence.ChannelRepositoryImpl{Postgres: s.app.Postgres}
}

func (s Setup) Cache() {
	s.app.Cache = cache.NewCacheImpl()
}

func (s Setup) Email() {
	s.app.Email = email.NewSesImpl()
}

func (s Setup) Logger(taskname string) {
	s.app.Logger, _ = logger.New(taskname)
}

func (s *Setup) Queue() {
	s.app.Queue = queue.NewKafkaImpl([]string{viper.GetString("KAFKA_BROKER")})
}

func (s Setup) Metrics() {
	s.app.Metrics = metrics.NewPrometheusImpl()
}

func (s Setup) Finish() {
	infra.App = *s.app
}
