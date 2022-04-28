package session

import (
	"github.com/heat1q/opt"

	"github.com/boardsite-io/server/internal/config"
	libErr "github.com/boardsite-io/server/pkg/errors"
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
		return libErr.ErrBadRequest.Wrap(libErr.WithErrorf("password cannot be longer than 64 characters"))
	}
	if numUsers, ok := c.MaxUsers.Some(); ok && (numUsers < 1 || numUsers > maxUsers) {
		return libErr.ErrBadRequest.Wrap(libErr.WithErrorf("incorrect maxUsers"))
	}
	return nil
}
