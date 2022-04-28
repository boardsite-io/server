package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/boardsite-io/server/internal/metrics"
)

func Metrics(metrics metrics.Handler) echo.MiddlewareFunc {
	return metrics.MiddlewareFunc()
}
