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
	"github.com/heat1q/boardsite/session/sessionfakes"
)

func TestSession(t *testing.T) {
	t.Run("extract scb and user", func(t *testing.T) {
		sessionId := "abcd1234"
		userId := "user1"

		scb := &sessionfakes.FakeController{}
		scb.IDReturns(sessionId)
		scb.GetUserReadyReturns(&types.User{
			ID: userId,
		}, nil)
		dispatcher := &sessionfakes.FakeDispatcher{}
		dispatcher.GetSCBReturns(scb, nil)

		e := echo.New()
		e.Use(middleware.Session(dispatcher))
		s := httptest.NewServer(e)
		defer s.Close()
		handler := func(c echo.Context) error {
			assert.Equal(t, sessionId, c.Get(session.SessionCtxKey).(session.Controller).ID())
			assert.Equal(t, userId, c.Get(session.UserCtxKey).(*types.User).ID)
			return c.NoContent(http.StatusOK)
		}
		e.GET("/b/:id", handler)

		r, _ := http.NewRequest(http.MethodGet, s.URL+"/b/"+sessionId, nil)
		r.Header.Set(middleware.HeaderUserID, userId)
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
		e.Use(middleware.Session(dispatcher))
		s := httptest.NewServer(e)
		defer s.Close()
		handler := func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}
		e.GET("/b/:id", handler)

		r, _ := http.NewRequest(http.MethodGet, s.URL+"/b/"+sessionId, nil)
		r.Header.Set(middleware.HeaderUserID, userId)
		resp, err := http.DefaultClient.Do(r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("wrong userId", func(t *testing.T) {
		sessionId := "abcd1234"
		userId := "user1"

		scb := &sessionfakes.FakeController{}
		scb.IDReturns(sessionId)
		scb.GetUserReadyReturns(nil, assert.AnError)
		dispatcher := &sessionfakes.FakeDispatcher{}
		dispatcher.GetSCBReturns(scb, nil)

		e := echo.New()
		e.Use(middleware.Session(dispatcher))
		s := httptest.NewServer(e)
		defer s.Close()
		handler := func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}
		e.GET("/b/:id", handler)

		r, _ := http.NewRequest(http.MethodGet, s.URL+"/b/"+sessionId, nil)
		r.Header.Set(middleware.HeaderUserID, userId)
		resp, err := http.DefaultClient.Do(r)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
