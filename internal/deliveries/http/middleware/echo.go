package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseWriter struct {
	Writer         io.Writer
	ResponseWriter http.ResponseWriter
}

func (r *ResponseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *ResponseWriter) WriteHeader(code int) {
	r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	return r.Writer.Write(b)
}

func (m *middleware) parseRequestBody(c echo.Context) []byte {
	var body []byte
	if c.Request().Body != nil {
		body, _ = io.ReadAll(c.Request().Body)
	}
	c.Request().Body = io.NopCloser(bytes.NewBuffer(body))
	return body
}

func (m *middleware) parseRequestHeader(c echo.Context) []byte {
	var header []byte
	if c.Request().Header != nil {
		header, _ = json.Marshal(c.Request().Header)
	}
	return header
}

func (m *middleware) getResponseBodyBuffer(c echo.Context) *bytes.Buffer {
	resBody := new(bytes.Buffer)
	mw := io.MultiWriter(c.Response().Writer, resBody)
	writer := &ResponseWriter{
		Writer:         mw,
		ResponseWriter: c.Response().Writer,
	}
	c.Response().Writer = writer
	return resBody
}
