package response

import (
	"github.com/labstack/echo/v4"
)

type CursorPaginationModel struct {
	Prev         string `json:"prev" example:"abc"`
	Next         string `json:"next" example:"cba"`
	TotalEntries int    `json:"totalEntries" example:"100"`
}

type PaginateContent[ModelOut any] interface {
	GetCursorTimestamp() string
	MappingData() ModelOut
}

func IsForward(c echo.Context) bool {
	return c.QueryParam("next_cursor") != ""
}

func IsBackward(c echo.Context) bool {
	return !IsForward(c) && c.QueryParam("prev_cursor") != ""
}

func NewCursorPagination[ModelOut any, S []E, E PaginateContent[ModelOut]](c echo.Context, collections S, hasMorePages bool, requestLimit, totalEntries int) CursorPaginationModel {
	var prevCursor, nextCursor string

	if len(collections) > 0 {
		if IsBackward(c) || hasMorePages {
			nextCursor = collections[len(collections)-1].GetCursorTimestamp()
		}

		if IsForward(c) || (hasMorePages && IsBackward(c)) {
			prevCursor = collections[0].GetCursorTimestamp()
		}
	}

	return CursorPaginationModel{
		Prev:         prevCursor,
		Next:         nextCursor,
		TotalEntries: totalEntries,
	}
}
