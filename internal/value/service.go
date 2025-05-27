package value

import (
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/spf13/viper"
)

const (
	// platforms

	EmailPlatform   = "email"
	SlackPlatform   = "slack"
	DiscordPlatform = "discord"

	// event status

	SuccessStatus = "success"
	ErrorStatus   = "error"
	PendingStatus = "pending"

	MaxRetries = 3
)

type NotificationInput struct {
	UUID     string   `json:"uuid" validate:"required"`
	Title    string   `json:"title" validate:"required"`
	Message  string   `json:"message"`
	Channels []string `json:"channels" validate:"required,min=1,max=20,dive,required"`
	Event    Event    `json:"event" validate:"required"`
}

type NotificationOutput struct {
	ID    int                 `json:"id"`
	UUID  string              `json:"uuid"`
	Body  entity.Notification `json:"body"`
	Error string              `json:"error"`
}

type Event struct {
	Name      string `json:"name" validate:"required"`
	Currency  string `json:"currency" validate:"required"`
	Requester string `json:"requester" validate:"required"`
	Receiver  string `json:"receiver" validate:"required"`
	Category  string `json:"category" validate:"required"`
	Timestamp int64  `json:"timestamp" validate:"required"`
	CostCents int64  `json:"cost_cents" validate:"required"`
}

func GetTopic() string {
	return viper.GetString("KAFKA_TOPIC")
}
