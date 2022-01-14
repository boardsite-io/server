package middleware

import (
	"github.com/heat1q/boardsite/session"
	"github.com/labstack/echo/v4"
)

func Session(dispatcher session.Dispatcher) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessionId := c.Param("id")
			if sessionId == "" {
				return echo.ErrNotFound
			}

			scb, err := dispatcher.GetSCB(sessionId)
			if err != nil {
				return echo.ErrForbidden
			}
			c.Set(session.SessionCtxKey, scb)

			userId := c.Request().Header.Get(HeaderUserID)
			if userId == "" {
				// userid could also be in params
				userId = c.Param("userId")
			}
			user, err := scb.GetUserReady(userId)
			if err != nil {
				return echo.ErrForbidden
			}
			c.Set(session.UserCtxKey, user)

			return next(c)
		}
	}
}
