package http

import (
	"fmt"
	"go-flip-life-style-products/internal/deliveries/http/balance"
	"go-flip-life-style-products/internal/deliveries/http/middleware"
	"go-flip-life-style-products/internal/deliveries/http/statements"
	"go-flip-life-style-products/internal/deliveries/http/transactions"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *httpApi) Route() {
	h.e.GET("/", func(c echo.Context) error {
		message := fmt.Sprintf("Welcome to %s", h.c.Cfg.App.Name)
		return c.String(http.StatusOK, message)
	})

	mid := middleware.NewMiddleware(h.c)
	h.e.Use(mid.Context())
	h.e.Use(mid.Logger())
	balance.Route(h.c, h.e)
	statements.Route(h.c, h.e)
	transactions.Route(h.c, h.e)
}
