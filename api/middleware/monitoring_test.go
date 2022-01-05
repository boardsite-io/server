package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/heat1q/boardsite/api/log"
	"github.com/heat1q/boardsite/api/middleware"
)

func TestMonitoring(t *testing.T) {
	t.Run("contains logger", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		hndl := func(c echo.Context) error {
			return c.NoContent(http.StatusNoContent)
		}
		fn := middleware.Monitoring()(hndl)
		err := fn(c)

		assert.NoError(t, err)

		_, ok := c.Request().Context().Value(log.ContextKey).(*zap.SugaredLogger)
		assert.True(t, ok)
	})
}
