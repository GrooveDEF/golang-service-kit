package echoserver

import (
	"github.com/definancialbr/golang-container-kit/pkg/container"
	"github.com/labstack/echo/v4"
)

func NewEchoServer(ctn *container.Container) *echo.Echo {
	server := echo.New()

	server.Logger.SetLevel(99)
	server.HideBanner = true

	server.Use(MetricsMiddleware(ctn.Metrics))
	server.Use(LoggingMiddleware(ctn.Logging))

	metricsHandler := ctn.Metrics.Handler()

	server.GET("/metrics", func(ctx echo.Context) error {
		metricsHandler.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})

	livenessHandler := ctn.Probes.LivenessHandler()

	server.GET("/healthz", func(ctx echo.Context) error {
		livenessHandler.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})

	readinessHandler := ctn.Probes.ReadinessHandler()

	server.GET("/healthz/ready", func(ctx echo.Context) error {
		readinessHandler.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})

	return server

}
