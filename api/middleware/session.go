package middleware

import (
	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/session"
	sessionHttp "github.com/heat1q/boardsite/session/http"
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
			c.Set(sessionHttp.SessionCtxKey, scb)

			userId := c.Request().Header.Get(types.HeaderUserID)
			user, ok := scb.GetUsers()[userId]
			if !ok {
				c.Error(apiErrors.ErrForbidden)
				return nil
			}
			c.Set(sessionHttp.UserCtxKey, user)

			c.Request().Header.Get(types.HeaderSessionSecret)
			c.Set(sessionHttp.SecretCtxKey, c.Request().Header.Get(types.HeaderSessionSecret))

			return next(c)
		}
	}
}

func Host() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			scb, ok := c.Get(sessionHttp.SessionCtxKey).(session.Controller)
			if !ok {
				c.Error(apiErrors.ErrForbidden)
				return nil
			}
			user, ok := c.Get(sessionHttp.UserCtxKey).(*session.User)
			if !ok {
				c.Error(apiErrors.ErrForbidden)
				return nil
			}
			secret := c.Get(sessionHttp.SecretCtxKey)

			if user.ID != scb.Config().Host || secret != scb.Config().Secret {
				c.Error(apiErrors.ErrForbidden)
				return nil
			}
			return next(c)
		}
	}
}
