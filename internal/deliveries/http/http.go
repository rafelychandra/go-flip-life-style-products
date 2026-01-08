package http

import (
	"context"
	"go-flip-life-style-products/internal/contract"
	"go-flip-life-style-products/internal/pkg/graceful"

	"github.com/labstack/echo/v4"
)

type (
	Http interface {
		Start() graceful.ProcessStarter
		Stop() graceful.ProcessStopper
	}
	httpApi struct {
		e *echo.Echo
		c *contract.Contract
	}
)

func New(c *contract.Contract) Http {
	return &httpApi{
		e: echo.New(),
		c: c,
	}
}

func (h *httpApi) Start() graceful.ProcessStarter {
	h.Route()
	return func() error {
		return h.e.Start(":" + h.c.Cfg.App.Port)
	}
}

func (h *httpApi) Stop() graceful.ProcessStopper {
	return graceful.ProcessStopper{
		Name: "HTTP",
		Stop: func(ctx context.Context) error {
			return h.e.Shutdown(ctx)
		},
	}
}
