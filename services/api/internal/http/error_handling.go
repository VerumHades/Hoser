package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var code int
		var message interface{}

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message
		} else {
			code = http.StatusInternalServerError
			message = "Internal Server Error"
		}

		c.JSON(code, map[string]interface{}{
			"error":  message,
			"status": code,
		})
	}
}
