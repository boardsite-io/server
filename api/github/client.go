package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/heat1q/boardsite/api/config"
	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/heat1q/boardsite/redis"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Client
type Client interface {
	PostToken(ctx context.Context, code string) (*TokenResponse, error)
	GetUserEmails(ctx context.Context, token string) ([]UserEmail, error)
}

type client struct {
	cfg   *config.Github
	cache redis.Handler
}

func NewClient(cfg *config.Github, cache redis.Handler) Client {
	return &client{
		cfg:   cfg,
		cache: cache,
	}
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (gh *client) PostToken(ctx context.Context, code string) (*TokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(echo.HeaderAccept, echo.MIMEApplicationJSON)

	reqQuery := url.Values{}
	reqQuery.Add("client_id", gh.cfg.ClientId)
	reqQuery.Add("client_secret", gh.cfg.ClientSecret)
	reqQuery.Add("code", code)
	req.URL.RawQuery = reqQuery.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, apiErrors.ErrBadGateway.Wrap(apiErrors.WithError(err))
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

type UserEmail struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

func (gh *client) GetUserEmails(ctx context.Context, token string) ([]UserEmail, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userEmailURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(echo.HeaderAuthorization, "bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, apiErrors.ErrBadGateway.Wrap(apiErrors.WithError(err))
	}
	defer resp.Body.Close()

	var emails []UserEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, err
	}

	return emails, nil
}
