package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gurodrigues-dev/notifier-app/internal/domain/repository"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/venture-technology/venture/pkg/stringcommon"
)

type Middleware struct {
	authRepository repository.AuthRepository
}

func NewMiddleware(
	authRepository repository.AuthRepository,
) *Middleware {
	return &Middleware{
		authRepository: authRepository,
	}
}

const (
	headerAdminToken = "Admin-Token"
	headerAuthToken  = "Authorization"
)

func (m *Middleware) AdminMiddleware() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		adminToken := strings.TrimSpace(httpContext.GetHeader(headerAdminToken))
		if stringcommon.Empty(adminToken) {
			httpContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "admin token is required"})
			return
		}

		if adminToken == viper.GetString("ADMIN_TOKEN") {
			httpContext.Next()
			httpContext.Set("isAdmin", true)
			return
		}

		httpContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid admin token"})
		return
	}
}

func (m *Middleware) TokenMiddleware() gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		token := strings.TrimSpace(httpContext.GetHeader(headerAuthToken))
		if stringcommon.Empty(token) {
			httpContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
			return
		}

		_, err := m.authRepository.GetToken(token)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpContext.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}

		if err != nil {
			httpContext.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		httpContext.Next()
		httpContext.Set("token", token)
		return
	}
}

func (m *Middleware) PrometheusMiddleware(metrics contracts.Metrics) gin.HandlerFunc {
	return func(httpContext *gin.Context) {
		start := time.Now()

		httpContext.Next()

		duration := time.Since(start).Seconds()
		method := httpContext.Request.Method
		path := httpContext.FullPath()
		metrics.IncRequest(method, path)
		metrics.ObserveDuration(method, path, duration)
	}
}
