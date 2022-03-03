package middleware

import (
	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/session"
	"github.com/heat1q/boardsite/session/http"
	"github.com/labstack/echo/v4"
)

func Session(dispatcher session.Dispatcher) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessionId := c.Param("id")
			if sessionId == "" {
				c.Error(apiErrors.ErrForbidden)
				return nil
			}

			scb, err := dispatcher.GetSCB(sessionId)
			if err != nil {
				c.Error(apiErrors.ErrNotFound)
				return nil
			}
			c.Set(http.SessionCtxKey, scb)

			userId := c.Request().Header.Get(types.HeaderUserID)
			if userId == "" {
				// userid could also be in params
				userId = c.Param("userId")
			}
			user, ok := scb.GetUsers()[userId]
			if !ok {
				c.Error(apiErrors.ErrForbidden)
				return nil
			}
			c.Set(http.UserCtxKey, user)

			c.Request().Header.Get(types.HeaderSessionSecret)
			c.Set(http.SecretCtxKey, c.Request().Header.Get(types.HeaderSessionSecret))

			return next(c)
		}
	}
}
