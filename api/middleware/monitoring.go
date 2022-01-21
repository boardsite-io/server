package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"

	"github.com/heat1q/boardsite/api/log"
)

func Monitoring(prom *prometheus.Prometheus) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			meta := make(map[string]interface{}, 2)

			meta[log.TagTraceID] = newTraceID()
			if sessionID := c.Param("id"); sessionID != "" {
				meta[log.TagSessionID] = sessionID
			}

			// set logger to request context
			ctx := log.WrapCtx(c.Request().Context(), meta)
			c.SetRequest(c.Request().WithContext(ctx))

			// collect prometheus metrics
			return prom.HandlerFunc(next)(c)
		}
	}
}

func newTraceID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Global().Warnf("generate new uuid failed: %v", err)
	}
	return id.String()
}
