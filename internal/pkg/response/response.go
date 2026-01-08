package response

import "github.com/labstack/echo/v4"

type (
	ErrorModel struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	RestPaginationResponseModel[T any] struct {
		Kind       string                 `json:"kind" example:"collection"`
		Contents   T                      `json:"contents"`
		Pagination *CursorPaginationModel `json:"pagination,omitempty"`
	}
)

func Success(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, data)
}

func Error(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, ErrorModel{
		Code:    statusCode,
		Message: message,
	})
}

func CursorPagination[ModelResponse any, S []E, E PaginateContent[ModelResponse]](c echo.Context, statusCode int, data S, requestLimit, totalRows int) error {
	hasMorePages := len(data) > (requestLimit - 1)

	if len(data) > 0 {
		if hasMorePages {
			data = data[:len(data)-1]
		}

		if IsBackward(c) {
			for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
				data[i], data[j] = data[j], data[i]
			}
		}
	}

	contents := make([]ModelResponse, 0)
	for _, d := range data {
		res := d.MappingData()
		if &res != nil {
			contents = append(contents, res)
		}
	}

	pagination := NewCursorPagination[ModelResponse](c, data, hasMorePages, requestLimit, totalRows)

	return c.JSON(statusCode, RestPaginationResponseModel[[]ModelResponse]{
		Kind:       "collection",
		Contents:   contents,
		Pagination: &pagination,
	})
}
