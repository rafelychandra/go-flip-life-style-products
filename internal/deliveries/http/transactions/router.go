package transactions

import (
	"go-flip-life-style-products/internal/contract"

	"github.com/labstack/echo/v4"
)

func Route(c *contract.Contract, e *echo.Echo) {
	group := e.Group("/transactions")
	NewTransactions(c).Route(group)
}
