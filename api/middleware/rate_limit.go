package middleware

import (
	"errors"
	"time"

	"github.com/heat1q/boardsite/api/types"

	apiErrors "github.com/heat1q/boardsite/api/errors"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type RateLimitingOption func(cfg *echomw.RateLimiterConfig)

func RateLimiting(rpm uint16, options ...RateLimitingOption) echo.MiddlewareFunc {
	if rpm == 0 {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return next(c)
			}
		}
	}
	memstoreCfg := echomw.RateLimiterMemoryStoreConfig{
		Rate:      rate.Limit(float64(rpm) / 60),
		Burst:     int(rpm),
		ExpiresIn: time.Minute,
	}
	cfg := echomw.RateLimiterConfig{
		Store:               echomw.NewRateLimiterMemoryStoreWithConfig(memstoreCfg),
		Skipper:             echomw.DefaultSkipper,
		IdentifierExtractor: ipExtractor,
		ErrorHandler: func(_ echo.Context, err error) error {
			return apiErrors.From(apiErrors.MissingIdentifier).Wrap(apiErrors.WithError(err))
		},
		DenyHandler: func(_ echo.Context, identifier string, err error) error {
			return apiErrors.From(apiErrors.RateLimitExceeded).Wrap(
				apiErrors.WithErrorf("rate limiter: denied %s: %w", identifier, err))
		},
	}

	for _, o := range options {
		o(&cfg)
	}

	return echomw.RateLimiterWithConfig(cfg)
}

// WithIP extracts the id from the real ip.
// This functional option is passed to RateLimiting.
func WithIP() RateLimitingOption {
	return func(cfg *echomw.RateLimiterConfig) {
		cfg.IdentifierExtractor = ipExtractor
	}
}

// WithUserID extracts the id from the userId header.
// This functional option is passed to RateLimiting.
func WithUserID() RateLimitingOption {
	return func(cfg *echomw.RateLimiterConfig) {
		cfg.IdentifierExtractor = userIDExtractor
	}
}

// WithUserIP extracts the id from the userId header + real ip.
// This functional option is passed to RateLimiting.
func WithUserIP() RateLimitingOption {
	return func(cfg *echomw.RateLimiterConfig) {
		fn := func(c echo.Context) (string, error) {
			userId, _ := userIDExtractor(c)
			ip, err := ipExtractor(c)
			if err != nil {
				return "", err
			}
			return userId + ":" + ip, nil
		}
		cfg.IdentifierExtractor = fn
	}
}

func ipExtractor(c echo.Context) (string, error) {
	return c.RealIP(), nil
}

func userIDExtractor(c echo.Context) (string, error) {
	userId := c.Request().Header.Get(types.HeaderUserID)
	if userId == "" {
		// userid could also be in params
		userId = c.Param("userId")
	}
	if userId == "" {
		return "", errors.New("no userId in header")
	}
	return userId, nil
}
