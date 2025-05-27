package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/gurodrigues-dev/notifier-app/cmd/api/notify-api/routes/v1"
	_ "github.com/gurodrigues-dev/notifier-app/docs"
	"github.com/gurodrigues-dev/notifier-app/internal/domain/middleware"
	"github.com/gurodrigues-dev/notifier-app/internal/infra"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/setup"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Notify API
// @version 0.1
// @description Notify API Swagger Documentation
// @termsOfService http://swagger.io/terms/

// @host localhost:9999
// @BasePath /api/v1/
func main() {
	setup := setup.NewSetup()
	setup.Logger("notify-app")
	setup.Cache()
	setup.Postgres()
	setup.Repositories()
	setup.Email()
	setup.Queue()
	setup.Metrics()

	setup.Finish()

	notifyApiPort := viper.GetString("NOTIFY_APP_PORT")
	notifyApi := setupNotifyApi()
	notifyApi.Run(fmt.Sprintf(":%s", notifyApiPort))
}

func setupNotifyApi() *gin.Engine {
	router := gin.Default()
	router.GET("/status", getStatus)
	router.GET("/metrics", infra.App.Metrics.ExposeHandler())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	groupV1 := router.Group("api/v1")
	groupV1.Use(configHeaders())

	middleware := middleware.NewMiddleware(infra.App.Repositories.AuthRepository)
	groupV1.Use(middleware.PrometheusMiddleware(infra.App.Metrics))

	v1.NewControllers().Routes(groupV1, middleware)
	return router
}

func getStatus(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("charset", "utf-8")
	c.Header("app_version", "2025.25.25 18:04")
	c.String(http.StatusOK, "ok, 25-05-25 version! =D")
}

func configHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, AdminToken")
		c.Header("Cross-Origin-Embedder-Policy", "require-corp")
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("charset", "utf-8")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
