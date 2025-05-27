package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gurodrigues-dev/notifier-app/internal/infra"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/setup"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/webhook"
	"github.com/gurodrigues-dev/notifier-app/internal/usecase"
	"github.com/spf13/viper"
)

func main() {
	setup := setup.NewSetup()
	setup.Logger("notify-dispatcher")
	setup.Cache()
	setup.Postgres()
	setup.Repositories()
	setup.Email()
	setup.Queue()
	setup.Metrics()

	setup.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		infra.App.Logger.Infof(fmt.Sprintf("Signal received: %s, shutting down gracefully...", sig))
		cancel()
	}()

	handler(ctx)

	<-ctx.Done()
	infra.App.Logger.Infof("Application stopped.")
}

func handler(ctx context.Context) {
	usecase := usecase.NewDispatcherUsecase(
		infra.App.Repositories.NotificationRepository,
		infra.App.Email,
		&webhook.DefaultClient{},
		infra.App.Logger,
	)

	infra.App.Logger.Infof("Starting Kafka dispatcher...")
	err := infra.App.Queue.Consumer(
		viper.GetString("KAFKA_TOPIC"),
		viper.GetString("KAFKA_GROUP"),
		func(message string) {
			err := usecase.Execute(message)
			if err != nil {
				infra.App.Logger.Errorf(fmt.Sprintf("Consume message error: %v", err))
			}
		},
	)
	if err != nil {
		infra.App.Logger.Errorf(fmt.Sprintf("Error starting consumer: %v", err))
	}
}
