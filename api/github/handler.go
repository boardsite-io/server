package github

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	apiErrors "github.com/heat1q/boardsite/api/errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/heat1q/boardsite/api/config"
	"github.com/heat1q/boardsite/redis"
)

const (
	AuthURL      = "https://github.com/login/oauth/authorize"
	tokenURL     = "https://github.com/login/oauth/access_token"
	apiURL       = "https://api.github.com/graphql"
	userEmailURL = "https://api.github.com/user/emails"
)

//counterfeiter:generate . Handler
type Handler interface {
	GetAuthorize(c echo.Context) error
	GetCallback(c echo.Context) error
}

type handler struct {
	cfg    *config.Configuration
	cache  redis.Handler
	client Client
}

func NewHandler(cfg *config.Configuration, cache redis.Handler, client Client) Handler {
	return &handler{
		cfg:    cfg,
		cache:  cache,
		client: client,
	}
}

// https://docs.github.com/en/developers/apps/building-oauth-apps/authorizing-oauth-apps#3-use-the-access-token-to-access-the-api
func (h *handler) GetAuthorize(c echo.Context) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	state := id.String()

	if err := h.cache.Put(c.Request().Context(), state, state, 10*time.Minute); err != nil {
		return fmt.Errorf("cache: put state: %w", err)
	}

	c.Response().Header().Set("Cache-Control", "no-store")

	query := url.Values{}
	query.Add("client_id", h.cfg.ClientId)
	query.Add("redirect_uri", fmt.Sprintf("%s:%d/github/oauth/callback", h.cfg.Server.BaseURL, h.cfg.Server.Port))
	query.Add("scope", strings.Join(h.cfg.Scope, " "))
	query.Add("state", state)
	return c.Redirect(http.StatusTemporaryRedirect, AuthURL+"?"+query.Encode())
}

func (h *handler) GetCallback(c echo.Context) error {
	var (
		ctx   = c.Request().Context()
		state = c.QueryParam("state")
		code  = c.QueryParam("code")
	)

	if s, err := h.cache.Get(ctx, state); err != nil || s == nil || string(s.([]byte)) != state {
		return apiErrors.ErrForbidden.Wrap(apiErrors.WithErrorf("compare states: %w", err))
	}
	_ = h.cache.Delete(ctx, state)

	tokenResp, err := h.client.PostToken(ctx, code)
	if err != nil {
		return err
	}

	query := url.Values{}
	query.Add("token_type", tokenResp.TokenType)
	query.Add("token", tokenResp.AccessToken)
	return c.Redirect(http.StatusTemporaryRedirect, h.cfg.Github.RedirectURI+"?"+query.Encode())
}
