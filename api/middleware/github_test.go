package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heat1q/boardsite/api/github"

	"github.com/heat1q/boardsite/api/github/githubfakes"
	"github.com/heat1q/boardsite/api/middleware"
	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/api/config"
	"github.com/labstack/echo/v4"
)

func TestGithubAuth(t *testing.T) {
	e := echo.New()
	e.HTTPErrorHandler = middleware.NewErrorHandler()
	validator := &githubfakes.FakeValidator{}
	tests := []struct {
		name        string
		cfg         *config.Github
		authHeader  string
		validateErr error
		wantErr     bool
	}{
		{
			name:       "successful validation",
			cfg:        &config.Github{Enabled: true},
			authHeader: "bearer abcd1234",
		},
		{
			name: "validation skipped",
			cfg:  &config.Github{Enabled: false},
		},
		{
			name:    "empty auth header",
			cfg:     &config.Github{Enabled: true},
			wantErr: true,
		},
		{
			name:       "malformed auth header",
			cfg:        &config.Github{Enabled: true},
			authHeader: "abcd1234",
			wantErr:    true,
		},
		{
			name:        "validation failed",
			cfg:         &config.Github{Enabled: true},
			authHeader:  "bearer abcd1234",
			validateErr: github.ErrNotValidated,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator.ValidateReturns(tt.validateErr)
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Add(echo.HeaderAuthorization, tt.authHeader)
			rr := httptest.NewRecorder()
			c := e.NewContext(r, rr)
			var hndlCalled bool
			hndl := func(c echo.Context) error {
				hndlCalled = true
				return c.NoContent(http.StatusNoContent)
			}

			fn := middleware.GithubAuth(tt.cfg, validator)(hndl)
			err := fn(c)

			if tt.wantErr {
				assert.False(t, hndlCalled)
			} else {
				assert.NoError(t, err)
				assert.True(t, hndlCalled)
			}
		})
	}
}
