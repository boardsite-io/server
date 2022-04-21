package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/boardsite-io/server/api/config"
	"github.com/boardsite-io/server/session"
	sessionHttp "github.com/boardsite-io/server/session/http"
	"github.com/boardsite-io/server/session/sessionfakes"
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

func Test_handler_PostCreateSessionConfig(t *testing.T) {
	const req = `
	{
		"config": {
			"maxUsers": 20,
			"readOnly": true,
			"password": "potato"
		}
	}`
	wantCfg := session.Config{
		Session: config.Session{
			MaxUsers: 20,
			ReadOnly: true,
		},
		Password: "potato",
	}
	scb := &sessionfakes.FakeController{}
	scb.ConfigReturns(wantCfg)
	dispatcher := &sessionfakes.FakeDispatcher{}
	dispatcher.CreateReturns(scb, nil)
	handler := sessionHttp.NewHandler(config.Session{}, dispatcher)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(req))
	rr := httptest.NewRecorder()
	c := echo.New().NewContext(r, rr)

	err := handler.PostCreateSessionConfig(c)

	assert.NoError(t, err)
	_, cfg := dispatcher.CreateArgsForCall(0)
	assert.Equal(t, wantCfg.MaxUsers, cfg.MaxUsers)
	assert.Equal(t, wantCfg.ReadOnly, cfg.ReadOnly)
	assert.Equal(t, wantCfg.Password, cfg.Password)
	var got session.CreateSessionResponse
	_ = json.NewDecoder(rr.Body).Decode(&got)
	assert.Equal(t, wantCfg, got.Config)
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
