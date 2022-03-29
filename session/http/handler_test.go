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
	cfg := session.Config{
		Session: config.Session{
			MaxUsers: 20,
			ReadOnly: true,
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
	scb := &sessionfakes.FakeController{}
	scb.ConfigReturns(cfg)
	dispatcher := &sessionfakes.FakeDispatcher{}
	dispatcher.CreateReturns(scb, nil)
	handler := sessionHttp.NewHandler(cfg.Session, dispatcher)
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"maxUsers": 20}`))
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rr := httptest.NewRecorder()
	c := e.NewContext(r, rr)
	c.Set(sessionHttp.SessionCtxKey, scb)

	err := handler.PutSessionConfig(c)

	assert.NoError(t, err)
}
