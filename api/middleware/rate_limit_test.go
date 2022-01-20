package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/api/middleware"
)

func getContext(e *echo.Echo) (*httptest.ResponseRecorder, echo.Context) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(echo.HeaderXForwardedFor, "127.0.0.1")
	r.Header.Set(middleware.HeaderUserID, "1234")
	rr := httptest.NewRecorder()
	return rr, e.NewContext(r, rr)
}

func handler(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func TestRateLimiting(t *testing.T) {
	e := echo.New()
	e.HTTPErrorHandler = middleware.NewErrorHandler()
	defer e.Close()

	t.Run("based on ip", func(t *testing.T) {
		fn := middleware.RateLimiting(1, middleware.WithIP())(handler)

		rr, c := getContext(e)
		_ = fn(c)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, c = getContext(e)
		_ = fn(c)
		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	})

	t.Run("based on usedId", func(t *testing.T) {
		fn := middleware.RateLimiting(1, middleware.WithUserID())(handler)

		rr, c := getContext(e)
		_ = fn(c)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, c = getContext(e)
		_ = fn(c)
		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	})

	t.Run("based on usedId plus ip", func(t *testing.T) {
		fn := middleware.RateLimiting(1, middleware.WithUserIP())(handler)

		rr, c := getContext(e)
		_ = fn(c)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, c = getContext(e)
		_ = fn(c)
		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	})
}
