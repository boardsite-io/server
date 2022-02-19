package github_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/heat1q/boardsite/api/github"
	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/api/config"
	"github.com/heat1q/boardsite/api/github/githubfakes"
	"github.com/heat1q/boardsite/redis/redisfakes"
	"github.com/labstack/echo/v4"
)

func Test_handler_GetAuthorize(t *testing.T) {
	e := echo.New()
	client := &githubfakes.FakeClient{}
	cache := &redisfakes.FakeHandler{}
	cfg := &config.Configuration{}
	cfg.Server.BaseURL = "http://localhost"
	cfg.Server.Port = 8000
	cfg.Github.ClientId = "client-Id"
	cfg.Github.Scope = []string{"user:email"}
	handler := github.NewHandler(cfg, cache, client)
	tests := []struct {
		name            string
		cachePutErr     error
		wantRedirectUrl string
		wantErr         bool
	}{
		{
			name:            "successful redirect to github auth",
			wantRedirectUrl: github.AuthURL + "?client_id=client-Id&redirect_uri=http%3A%2F%2Flocalhost%3A8000%2Fgithub%2Foauth%2Fcallback&scope=user%3Aemail",
		},
		{
			name:        "cache error",
			cachePutErr: assert.AnError,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.PutReturns(tt.cachePutErr)
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			c := e.NewContext(r, rr)

			err := handler.GetAuthorize(c)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, rr.Header().Get(echo.HeaderLocation), tt.wantRedirectUrl)
			}
		})
	}
}

func Test_handler_GetCallback(t *testing.T) {
	const mockState = "1234"
	e := echo.New()
	client := &githubfakes.FakeClient{}
	cache := &redisfakes.FakeHandler{}
	cfg := &config.Configuration{}
	cfg.Github.RedirectURI = "https://boardsite.io"
	handler := github.NewHandler(cfg, cache, client)
	tests := []struct {
		name            string
		cacheGet        interface{}
		cacheGetErr     error
		tokenResp       github.TokenResponse
		tokenErr        error
		wantRedirectUrl string
		wantErr         bool
	}{
		{
			name:     "successful redirect",
			cacheGet: []byte(mockState),
			tokenResp: github.TokenResponse{
				TokenType:   "bearer",
				AccessToken: "abcd1234",
			},
			wantRedirectUrl: "https://boardsite.io?token=abcd1234&token_type=bearer",
		},
		{
			name:     "mismatching state returns error",
			cacheGet: []byte("5678"),
			wantErr:  true,
		},
		{
			name:     "failing token response returns error",
			cacheGet: []byte(mockState),
			tokenErr: assert.AnError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.GetReturns(tt.cacheGet, tt.cacheGetErr)
			client.PostTokenReturns(&tt.tokenResp, tt.tokenErr)
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			q := url.Values{}
			q.Add("state", mockState)
			r.URL.RawQuery = q.Encode()
			rr := httptest.NewRecorder()
			c := e.NewContext(r, rr)

			err := handler.GetCallback(c)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, rr.Header().Get(echo.HeaderLocation), tt.wantRedirectUrl)
			}
		})
	}
}
