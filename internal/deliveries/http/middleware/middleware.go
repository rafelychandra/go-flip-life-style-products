package middleware

import (
	"go-flip-life-style-products/internal/contract"

	"github.com/labstack/echo/v4"
)

type (
	middleware struct {
		contract *contract.Contract
	}

	Middleware interface {
		Logger() echo.MiddlewareFunc
		Context() echo.MiddlewareFunc
	}
)

func NewMiddleware(contract *contract.Contract) Middleware {
	return &middleware{
		contract: contract,
	}
}
