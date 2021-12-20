package log

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const (
	ContextKey = "ctx-logger"
)

var (
	currentLogger = New()
)

func New() echo.Logger {
	logger := log.New("")
	logger.SetHeader("${time_rfc3339} ${level} ⇨")
	return logger
}

func Ctx(ctx context.Context) echo.Logger {
	logger, ok := ctx.Value(ContextKey).(echo.Logger)
	if !ok {
		currentLogger.Warn("Ctx doesn't contain logger")
		return currentLogger
	}
	return logger
}
