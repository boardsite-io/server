package github_test

import (
	"context"
	"testing"

	"github.com/gomodule/redigo/redis"

	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/api/github/githubfakes"
	"github.com/heat1q/boardsite/redis/redisfakes"

	"github.com/heat1q/boardsite/api/config"

	"github.com/heat1q/boardsite/api/github"
)

func Test_validator_Validate(t *testing.T) {
	ctx := context.Background()
	const mockToken = "abcd1234"
	client := &githubfakes.FakeClient{}
	cache := &redisfakes.FakeHandler{}
	cfg := &config.Github{
		Enabled: true,
		WhitelistedEmails: map[string]struct{}{
			"potato@boardsite.io": {},
		},
	}
	validator := github.NewValidator(cfg, cache, client)
	tests := []struct {
		name        string
		respEmails  []github.UserEmail
		respErr     error
		cacheGet    interface{}
		cacheGetErr error
		cachePutErr error
		wantErr     bool
	}{
		{
			name: "successful validation via emails",
			respEmails: []github.UserEmail{
				{
					Email:    "potato@boardsite.io",
					Verified: true,
				},
			},
			cacheGetErr: redis.ErrNil,
		},
		{
			name:     "successful validation due to cached token",
			cacheGet: []byte(mockToken),
		},
		{
			name: "fails with unverified email",
			respEmails: []github.UserEmail{
				{
					Email: "potato@boardsite.io",
				},
			},
			cacheGetErr: redis.ErrNil,
			wantErr:     true,
		},
		{
			name: "fails with unverified cache error",
			respEmails: []github.UserEmail{
				{
					Email:    "potato@boardsite.io",
					Verified: true,
				},
			},
			cacheGetErr: redis.ErrNil,
			cachePutErr: assert.AnError,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.GetUserEmailsReturns(tt.respEmails, tt.respErr)
			cache.GetReturns(tt.cacheGet, tt.cacheGetErr)
			cache.PutReturns(tt.cachePutErr)
			err := validator.Validate(ctx, mockToken)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
