package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/heat1q/boardsite/api/config"
	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/heat1q/boardsite/api/github"
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
				return apiErrors.ErrUnauthorized.Wrap(apiErrors.WithError(err))
			}

			return next(c)
		}
	}
}
