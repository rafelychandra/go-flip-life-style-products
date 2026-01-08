package balance

import (
	"go-flip-life-style-products/internal/contract"
	"go-flip-life-style-products/internal/models"
	"go-flip-life-style-products/internal/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Balance struct {
	c *contract.Contract
}

func NewBalance(c *contract.Contract) *Balance {
	return &Balance{
		c: c,
	}
}

func (b *Balance) Route(g *echo.Group) {
	g.GET("", b.getBalanceByUploadID)
}

func (b *Balance) getBalanceByUploadID(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		payload models.Balance
	)

	err := c.Bind(&payload)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	balance, err := b.c.Service.Balance.GetBalanceByUploadID(ctx, payload.UploadID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, balance)
}
