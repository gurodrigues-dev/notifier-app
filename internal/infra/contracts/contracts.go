package contracts

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type SESIface interface {
	SendEmail(email *entity.Email) error
	VerifyEmail(email string) error
}

type PostgresIface interface {
	Client() *gorm.DB
	Close() error
}

type Cacher interface {
	Set(key string, value any, expiration time.Duration) error
	Get(key string) (string, error)
	Expire(key string, expiration time.Duration) (bool, error)
}

type Logger interface {
	Infof(format string, args ...zap.Field)
	Errorf(format string, args ...zap.Field)
}

type Queue interface {
	Produce(topic, message string) error
	Consumer(topic, group string, handler func(message string)) error
}

type Metrics interface {
	IncRequest(method, path string)
	IncError(method, path string, errType string)
	ObserveDuration(method, path string, durationSeconds float64)
	ExposeHandler() gin.HandlerFunc
}
