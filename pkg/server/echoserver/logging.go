package echoserver

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/GrooveDEF/golang-container-kit/pkg/logging"
)

func LoggingMiddleware(log logging.LoggingService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()

			err := next(ctx)
			if err != nil {
				ctx.Error(err)
			}

			stop := time.Now()

			req := ctx.Request()
			res := ctx.Response()

			fields := []interface{}{
				"remote_ip", ctx.RealIP(),
				"host", req.Host,
				"status", res.Status,
				"user_agent", req.UserAgent(),
				"uri", req.RequestURI,
				"method", req.Method,
				"path", req.URL.Path,
				"referer", req.Referer(),
				"latency", strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
				"latency_human", stop.Sub(start).String(),
				"bytes_in", req.Header.Get(echo.HeaderContentLength),
				"bytes_out", strconv.FormatInt(res.Size, 10),
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id != "" {
				fields = append(fields, "request_id", id)
			}

			n := res.Status
			switch {
			case n >= 500:
				log.Error(err.Error(), fields...)
			case n >= 400:
				log.Warn(err.Error(), fields...)
			case n >= 300:
				log.Error("Redirection", fields...)
			default:
				log.Debug("Success", fields...)
			}

			return nil
		}
	}
}
