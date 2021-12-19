package middleware

import (
	"github.com/heat1q/boardsite/api/errors"
	"github.com/labstack/echo/v4"
)

func GetCustomHTTPErrorHandler(echoServer *echo.Echo) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		if e, ok := err.(*errors.HTTPError); ok {
			if e.Message != "" {
				err = echo.NewHTTPError(e.Status, e.Message)
			} else {
				err = echo.NewHTTPError(e.Status)
			}
		}
		echoServer.DefaultHTTPErrorHandler(err, c)
	}
}
