package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/api/middleware"
	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/session"
	sessionHttp "github.com/heat1q/boardsite/session/http"
	"github.com/heat1q/boardsite/session/sessionfakes"
)

func TestSession(t *testing.T) {
	t.Run("extract scb and user", func(t *testing.T) {
		sessionId := "abcd1234"
		userId := "user1"

		scb := &sessionfakes.FakeController{}
		scb.IDReturns(sessionId)
		scb.GetUsersReturns(map[string]*session.User{
			userId: {ID: userId},
		})
		scb.AllowReturns(true)
		dispatcher := &sessionfakes.FakeDispatcher{}
		dispatcher.GetSCBReturns(scb, nil)

		e := echo.New()
		e.HTTPErrorHandler = middleware.NewErrorHandler()
		e.Use(middleware.Session(dispatcher))
		s := httptest.NewServer(e)
		defer s.Close()
		handler := func(c echo.Context) error {
			assert.Equal(t, sessionId, c.Get(sessionHttp.SessionCtxKey).(session.Controller).ID())
			assert.Equal(t, userId, c.Get(sessionHttp.UserCtxKey).(*session.User).ID)
			return c.NoContent(http.StatusOK)
		}
		e.GET("/b/:id", handler)

		r, _ := http.NewRequest(http.MethodGet, s.URL+"/b/"+sessionId, nil)
		r.Header.Set(types.HeaderUserID, userId)
		resp, err := http.DefaultClient.Do(r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("wrong sessionId", func(t *testing.T) {
		sessionId := "abcd1234"
		userId := "user1"

		dispatcher := &sessionfakes.FakeDispatcher{}
		dispatcher.GetSCBReturns(nil, assert.AnError)

		e := echo.New()
		e.HTTPErrorHandler = middleware.NewErrorHandler()
		e.Use(middleware.Session(dispatcher))
		s := httptest.NewServer(e)
		defer s.Close()
		handler := func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}
		e.GET("/b/:id", handler)

		r, _ := http.NewRequest(http.MethodGet, s.URL+"/b/"+sessionId, nil)
		r.Header.Set(types.HeaderUserID, userId)
		resp, err := http.DefaultClient.Do(r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("wrong userId", func(t *testing.T) {
		sessionId := "abcd1234"
		userId := "user1"

		scb := &sessionfakes.FakeController{}
		scb.IDReturns(sessionId)
		scb.GetUsersReturns(map[string]*session.User{})
		scb.AllowReturns(true)
		dispatcher := &sessionfakes.FakeDispatcher{}
		dispatcher.GetSCBReturns(scb, nil)

		e := echo.New()
		e.HTTPErrorHandler = middleware.NewErrorHandler()
		e.Use(middleware.Session(dispatcher))
		s := httptest.NewServer(e)
		defer s.Close()
		handler := func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}
		e.GET("/b/:id", handler)

		r, _ := http.NewRequest(http.MethodGet, s.URL+"/b/"+sessionId, nil)
		r.Header.Set(types.HeaderUserID, userId)
		resp, err := http.DefaultClient.Do(r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func TestHost(t *testing.T) {
	e := echo.New()
	cfg := session.Config{
		Host:   "userId",
		Secret: "secret",
	}
	scb := &sessionfakes.FakeController{}
	scb.ConfigReturns(cfg)
	tests := []struct {
		name    string
		store   map[string]interface{}
		wantErr bool
	}{
		{
			name: "happy path",
			store: map[string]interface{}{
				sessionHttp.SessionCtxKey: scb,
				sessionHttp.UserCtxKey:    &session.User{ID: "userId"},
				sessionHttp.SecretCtxKey:  "secret",
			},
		},
		{
			name:    "missing context",
			store:   map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "missing user",
			store: map[string]interface{}{
				sessionHttp.SessionCtxKey: scb,
				sessionHttp.SecretCtxKey:  "secret",
			},
			wantErr: true,
		},
		{
			name: "wrong user",
			store: map[string]interface{}{
				sessionHttp.SessionCtxKey: scb,
				sessionHttp.UserCtxKey:    &session.User{ID: "userId1234"},
				sessionHttp.SecretCtxKey:  "secret",
			},
			wantErr: true,
		},
		{
			name: "wrong secret",
			store: map[string]interface{}{
				sessionHttp.SessionCtxKey: scb,
				sessionHttp.UserCtxKey:    &session.User{ID: "userId"},
				sessionHttp.SecretCtxKey:  "1234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			c := e.NewContext(r, rr)
			for k, v := range tt.store {
				c.Set(k, v)
			}
			handlerCalled := false
			handler := func(c echo.Context) error {
				handlerCalled = true
				return c.NoContent(http.StatusOK)
			}

			err := middleware.Host()(handler)(c)

			assert.NoError(t, err)
			if tt.wantErr {
				assert.False(t, handlerCalled)
			} else {
				assert.True(t, handlerCalled)
			}
		})
	}
}
