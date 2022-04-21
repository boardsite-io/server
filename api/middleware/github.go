package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/boardsite-io/server/api/config"
	apiErrors "github.com/boardsite-io/server/api/errors"
	"github.com/boardsite-io/server/api/github"
)

func GithubAuth(cfg *config.Github, validator github.Validator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !cfg.Enabled {
				return next(c)
			}

			h := c.Request().Header.Get(echo.HeaderAuthorization)
			auth := strings.Split(h, " ")
			if len(auth) != 2 {
				return apiErrors.ErrForbidden
			}
			token := auth[1]

			if err := validator.Validate(c.Request().Context(), token); err != nil {
				c.Error(apiErrors.ErrUnauthorized.Wrap(apiErrors.WithError(err)))
				return nil
			}

			return next(c)
		}
	}
}
