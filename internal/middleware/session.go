package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/boardsite-io/server/internal/session"
	sessionHttp "github.com/boardsite-io/server/internal/session/http"
	"github.com/boardsite-io/server/pkg/constant"
	libErr "github.com/boardsite-io/server/pkg/errors"
)

func Session(dispatcher session.Dispatcher) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessionId := c.Param("id")
			if sessionId == "" {
				c.Error(libErr.ErrForbidden)
				return nil
			}

			scb, err := dispatcher.GetSCB(sessionId)
			if err != nil {
				c.Error(libErr.ErrNotFound)
				return nil
			}
			c.Set(sessionHttp.SessionCtxKey, scb)

			userId := c.Request().Header.Get(constant.HeaderUserID)
			user, ok := scb.GetUsers()[userId]
			if !ok {
				c.Error(libErr.ErrForbidden)
				return nil
			}
			c.Set(sessionHttp.UserCtxKey, user)
			c.Set(sessionHttp.SecretCtxKey, c.Request().Header.Get(constant.HeaderSessionSecret))

			if !sessionHttp.AllowUser(c) {
				c.Error(libErr.ErrForbidden)
				return nil
			}

			return next(c)
		}
	}
}

func Host() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			scb, ok := c.Get(sessionHttp.SessionCtxKey).(session.Controller)
			if !ok {
				c.Error(libErr.ErrForbidden)
				return nil
			}
			user, ok := c.Get(sessionHttp.UserCtxKey).(*session.User)
			if !ok {
				c.Error(libErr.ErrForbidden)
				return nil
			}
			secret := c.Get(sessionHttp.SecretCtxKey)

			if user.ID != scb.Config().Host || secret != scb.Config().Secret {
				c.Error(libErr.ErrForbidden)
				return nil
			}
			return next(c)
		}
	}
}
