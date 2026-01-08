package middleware

import (
	"fmt"
	"time"

	loggerPkg "go-flip-life-style-products/internal/pkg/logger"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (m *middleware) Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			ctx := c.Request().Context()

			req := c.Request()
			res := c.Response()

			reqBody := m.parseRequestBody(c)
			reqHeader := m.parseRequestHeader(c)
			resBodyBuff := m.getResponseBodyBuffer(c)

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			latency := time.Since(start)

			fields := log.Fields{
				"start_time":     start.Format(time.RFC3339Nano),
				"end_time":       start.Add(latency).Format(time.RFC3339Nano),
				"method":         req.Method,
				"url_path":       req.URL.String(),
				"request_body":   string(reqBody),
				"request_header": string(reqHeader),
				"status":         res.Status,
				"response":       string(resBodyBuff.Bytes()),
				"latency":        latency.String(),
			}

			message := fmt.Sprintf(
				"%d %s %s %s",
				res.Status,
				req.Method,
				req.URL.String(),
				latency,
			)

			switch {
			case res.Status >= 500:
				loggerPkg.Error(ctx, message, fields)
			case res.Status >= 400:
				loggerPkg.Warn(ctx, message, fields)
			case res.Status >= 300:
				loggerPkg.Warn(ctx, message, fields)
			default:
				loggerPkg.Info(ctx, message, fields)
			}

			return nil
		}
	}
}
