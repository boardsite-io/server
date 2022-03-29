package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/heat1q/boardsite/api/log"
	"github.com/heat1q/boardsite/api/metrics"
)

func Monitoring() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			meta := make(map[string]any, 2)

			meta[log.TagTraceID] = newTraceID()
			if sessionID := c.Param("id"); sessionID != "" {
				meta[log.TagSessionID] = sessionID
			}

			// set logger to request context
			ctx := log.WrapCtx(c.Request().Context(), meta)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func Metrics(metrics metrics.Handler) echo.MiddlewareFunc {
	return metrics.MiddlewareFunc()
}

func newTraceID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Global().Warnf("generate new uuid failed: %v", err)
	}
	return id.String()
}
