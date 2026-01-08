package transactions

import (
	"go-flip-life-style-products/internal/contract"
	"go-flip-life-style-products/internal/models"
	"go-flip-life-style-products/internal/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Transactions struct {
	c *contract.Contract
}

func NewTransactions(c *contract.Contract) *Transactions {
	return &Transactions{
		c: c,
	}
}

func (t *Transactions) Route(g *echo.Group) {
	g.GET("/issues", t.getListTransactionsIssues)
}

func (t *Transactions) getListTransactionsIssues(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		payload models.ReqGetListIssuesTransaction
	)

	err := c.Bind(&payload)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	resp, count, limit, err := t.c.Service.Transaction.GetListIssuesTransaction(ctx, payload)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.CursorPagination[models.Transaction](c, http.StatusOK, resp, limit, count)
}
