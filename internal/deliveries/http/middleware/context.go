package middleware

import (
	"context"
	uuidPkg "go-flip-life-style-products/internal/pkg/uuid"

	"github.com/labstack/echo/v4"
)

func (m *middleware) Context() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				ctx = c.Request().Context()
			)

			ctx = context.WithValue(ctx, uuidPkg.CorrelationIDKey, uuidPkg.UUID())
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
