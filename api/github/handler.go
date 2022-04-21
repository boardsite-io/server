package github

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/boardsite-io/server/api/config"
	"github.com/boardsite-io/server/redis"
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

	query := url.Values{}
	query.Add("client_id", h.cfg.ClientId)
	query.Add("redirect_uri", fmt.Sprintf("%s/github/oauth/callback", h.cfg.Server.BaseURL))
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

	if c.QueryParam("error") != "" {
		return h.queryError(c, c.QueryParam("error_description"))
	}

	if s, err := h.cache.Get(ctx, state); err != nil || s == nil || string(s.([]byte)) != state {
		return h.queryError(c, "mismatching states")
	}
	_ = h.cache.Delete(ctx, state)

	tokenResp, err := h.client.PostToken(ctx, code)
	if err != nil {
		return h.queryError(c, "failed to post token")
	}

	query := url.Values{}
	query.Add("token_type", tokenResp.TokenType)
	query.Add("token", tokenResp.AccessToken)
	return c.Redirect(http.StatusTemporaryRedirect, h.cfg.Github.RedirectURI+"?"+query.Encode())
}

func (h *handler) queryError(c echo.Context, message string) error {
	query := url.Values{}
	query.Add("error", message)
	return c.Redirect(http.StatusTemporaryRedirect, h.cfg.Github.RedirectURI+"?"+query.Encode())
}
