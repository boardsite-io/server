package http

import (
	"net/http"

	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/labstack/echo/v4"

	"github.com/heat1q/boardsite/session"
)

const (
	SessionCtxKey = "boardsite-session"
	SecretCtxKey  = "boardsite-session-secret"
	UserCtxKey    = "boardsite-user"
)

func AllowUser(c echo.Context) bool {
	scb, err := getSCB(c)
	if err != nil {
		return false
	}
	user, ok := c.Get(UserCtxKey).(*session.User)
	if !ok {
		return false
	}
	secret, _ := c.Get(SecretCtxKey).(string)

	// request additionally need to check for correct secret
	if scb.Config().ReadOnly && c.Request().Method != http.MethodGet &&
		(user.ID != scb.Config().Host || secret != scb.Config().Secret) {
		return false
	}

	return scb.Allow(scb.Config().Host) // only check pw
}

func getSCB(c echo.Context) (session.Controller, error) {
	scb, ok := c.Get(SessionCtxKey).(session.Controller)
	if !ok {
		return nil, apiErrors.ErrForbidden
	}
	return scb, nil
}

func getUser(c echo.Context) (*session.User, error) {
	u, ok := c.Get(UserCtxKey).(*session.User)
	if !ok {
		return nil, apiErrors.ErrForbidden
	}
	return u, nil
}
