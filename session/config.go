package session

import (
	"github.com/heat1q/opt"

	"github.com/heat1q/boardsite/api/config"
	apiErrors "github.com/heat1q/boardsite/api/errors"
)

// TODO move to config
const maxUsers = 50

type Config struct {
	ID     string `json:"id"`
	Host   string `json:"host,omitempty"`
	Secret string `json:"-"`

	config.Session
	Password string `json:"password"`
}

func (c *Config) Update(incoming *ConfigRequest) error {
	if incoming == nil {
		return nil
	}
	if err := incoming.Validate(); err != nil {
		return err
	}
	if numUsers, ok := incoming.MaxUsers.Some(); ok {
		c.MaxUsers = numUsers
	}
	if ro, ok := incoming.ReadOnly.Some(); ok {
		c.ReadOnly = ro
	}
	if pw, ok := incoming.Password.Some(); ok {
		c.Password = pw
	}
	return nil
}

type ConfigRequest struct {
	MaxUsers opt.Option[int]    `json:"maxUsers,omitempty"`
	ReadOnly opt.Option[bool]   `json:"readOnly,omitempty"`
	Password opt.Option[string] `json:"password,omitempty"`
}

func (c *ConfigRequest) Validate() error {
	if pw, ok := c.Password.Some(); ok && len(pw) > 2<<5 {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("password cannot be longer than 64 characters"))
	}
	if numUsers, ok := c.MaxUsers.Some(); ok && (numUsers < 1 || numUsers > maxUsers) {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("incorrect maxUsers"))
	}
	return nil
}
