package github

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/boardsite-io/server/api/config"
	"github.com/boardsite-io/server/redis"
)

var ErrNotValidated = errors.New("validator: not validated")

//counterfeiter:generate . Validator
type Validator interface {
	Validate(ctx context.Context, token string) error
}

type validator struct {
	cfg    *config.Github
	cache  redis.Handler
	client Client
}

func NewValidator(cfg *config.Github, cache redis.Handler, client Client) Validator {
	return &validator{
		cfg:    cfg,
		cache:  cache,
		client: client,
	}
}

func (v *validator) Validate(ctx context.Context, token string) error {
	if t, err := v.cache.Get(ctx, token); err == nil && t != nil && string(t.([]byte)) == token {
		return nil
	}

	if err := v.userEmail(ctx, token); err == nil {
		return v.cache.Put(ctx, token, token, 24*time.Hour)
	}

	if err := v.userSponsor(ctx, token); err == nil {
		return v.cache.Put(ctx, token, token, time.Hour)
	}

	return ErrNotValidated
}

func (v *validator) userEmail(ctx context.Context, token string) error {
	emails, err := v.client.GetUserEmails(ctx, token)
	if err != nil {
		return fmt.Errorf("validator: getuseremails: %w", err)
	}

	for _, e := range emails {
		if !e.Verified {
			continue
		}
		if _, ok := v.cfg.WhitelistedEmails[e.Email]; ok {
			return nil
		}
	}

	return ErrNotValidated
}

func (v *validator) userSponsor(ctx context.Context, token string) error {
	return ErrNotValidated
}
