package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/api/config"
	"github.com/heat1q/boardsite/session"
	sessionHttp "github.com/heat1q/boardsite/session/http"
	"github.com/heat1q/boardsite/session/sessionfakes"
)

func Test_handler_PostCreateSession(t *testing.T) {
	readonly := true
	cfg := session.Config{
		Session: config.Session{
			MaxUsers: 20,
			ReadOnly: &readonly,
		},
	}
	scb := &sessionfakes.FakeController{}
	scb.ConfigReturns(cfg)
	dispatcher := &sessionfakes.FakeDispatcher{}
	dispatcher.CreateReturns(scb, nil)
	handler := sessionHttp.NewHandler(cfg.Session, dispatcher)

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()
	c := echo.New().NewContext(r, rr)

	err := handler.PostCreateSession(c)

	assert.NoError(t, err)
	var got session.CreateSessionResponse
	_ = json.NewDecoder(rr.Body).Decode(&got)
	assert.Equal(t, cfg, got.Config)
}

func Test_handler_PutSessionConfig(t *testing.T) {
	e := echo.New()
	cfg := session.Config{
		Host:   "userId",
		Secret: "secret",
	}
	user := session.User{
		ID:    "userId",
		Alias: "test",
		Color: "#00ff00",
	}
	scb := &sessionfakes.FakeController{}
	scb.ConfigReturns(cfg)
	dispatcher := &sessionfakes.FakeDispatcher{}
	dispatcher.CreateReturns(scb, nil)
	handler := sessionHttp.NewHandler(cfg.Session, dispatcher)

	t.Run("successful", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"maxUsers": 20}`))
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		rr := httptest.NewRecorder()
		c := e.NewContext(r, rr)
		c.Set(sessionHttp.SessionCtxKey, scb)
		c.Set(sessionHttp.UserCtxKey, &user)
		c.Set(sessionHttp.SecretCtxKey, "secret")

		err := handler.PutSessionConfig(c)

		assert.NoError(t, err)
	})

	t.Run("missing/wrong secret", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"maxUsers": 20}`))
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		rr := httptest.NewRecorder()
		c := e.NewContext(r, rr)
		c.Set(sessionHttp.SessionCtxKey, scb)
		c.Set(sessionHttp.UserCtxKey, &user)
		c.Set(sessionHttp.SecretCtxKey, "1234")

		err := handler.PutSessionConfig(c)

		assert.Error(t, err)
	})

	t.Run("user not host", func(t *testing.T) {
		user := session.User{
			ID:    "userId2",
			Alias: "test",
			Color: "#00ff00",
		}
		r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"maxUsers": 20}`))
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		rr := httptest.NewRecorder()
		c := e.NewContext(r, rr)
		c.Set(sessionHttp.SessionCtxKey, scb)
		c.Set(sessionHttp.UserCtxKey, &user)
		c.Set(sessionHttp.SecretCtxKey, "secret")

		err := handler.PutSessionConfig(c)

		assert.Error(t, err)
	})
}
