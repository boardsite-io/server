package middleware

import (
	"errors"
	"net/http"

	apiErrors "github.com/boardsite-io/server/api/errors"
	"github.com/boardsite-io/server/api/log"
	"github.com/labstack/echo/v4"
)

func NewErrorHandler() func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		httpErr := &apiErrors.HTTPError{}
		if ok := errors.As(err, &httpErr); !ok {
			echoErr := &echo.HTTPError{}
			if ok := errors.As(err, &echoErr); !ok {
				httpErr = apiErrors.ErrInternalServerError
			} else {
				httpErr.Status = echoErr.Code
			}
		}

		if httpErr.Message == "" {
			httpErr.Message = http.StatusText(httpErr.Status)
		}

		if err := c.JSON(httpErr.Status, httpErr); err != nil {
			log.Ctx(c.Request().Context()).Warn("failed to write error response")
		}
	}
}
