package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	http2 "github.com/heat1q/boardsite/session/http"

	"github.com/heat1q/boardsite/api/types"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/api/middleware"
	"github.com/heat1q/boardsite/session"
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
		dispatcher := &sessionfakes.FakeDispatcher{}
		dispatcher.GetSCBReturns(scb, nil)

		e := echo.New()
		e.HTTPErrorHandler = middleware.NewErrorHandler()
		e.Use(middleware.Session(dispatcher))
		s := httptest.NewServer(e)
		defer s.Close()
		handler := func(c echo.Context) error {
			assert.Equal(t, sessionId, c.Get(http2.SessionCtxKey).(session.Controller).ID())
			assert.Equal(t, userId, c.Get(http2.UserCtxKey).(*session.User).ID)
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
