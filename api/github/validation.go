package github

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/heat1q/boardsite/api/config"
	"github.com/heat1q/boardsite/redis"
)

type Validator interface {
	UserEmail(ctx context.Context, token string) error
}

type validator struct {
	cfg    *config.Github
	cache  redis.Handler
	client Client
}

func NewValidator(cfg *config.Github, cache redis.Handler) Validator {
	return &validator{
		cfg:    cfg,
		cache:  cache,
		client: NewClient(cfg, cache),
	}
}

func (v *validator) UserEmail(ctx context.Context, token string) error {
	if _, err := v.cache.Get(ctx, token); err == nil {
		return nil
	}

	emails, err := v.client.GetUserEmails(ctx, token)
	if err != nil {
		return fmt.Errorf("validator: getuseremails: %w", err)
	}

	if err := v.cache.Put(ctx, token, token, 24*time.Hour); err != nil {
		return fmt.Errorf("validator: put token cache: %w", err)
	}

	for _, e := range emails {
		if !e.Verified {
			continue
		}
		if _, ok := v.cfg.WhitelistedEmails[e.Email]; ok {
			return nil
		}
	}

	return errors.New("validator: email not validated")
}
