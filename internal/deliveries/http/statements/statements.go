package statements

import (
	"go-flip-life-style-products/internal/contract"
	"go-flip-life-style-products/internal/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Statements struct {
	c *contract.Contract
}

func NewStatements(c *contract.Contract) *Statements {
	return &Statements{
		c: c,
	}
}

func (s *Statements) Route(g *echo.Group) {
	g.POST("", s.upload)
}

func (s *Statements) upload(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	uploadID, err := s.c.Service.Statements.Upload(ctx, file)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, map[string]interface{}{
		"upload_id": uploadID,
	})
}
