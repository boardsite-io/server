package session

import (
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
	if incoming.ReadOnly != nil {
		c.ReadOnly = *incoming.ReadOnly
	}
	if incoming.Password != nil {
		c.Password = *incoming.Password
	}
	return nil
}

type ConfigRequest struct {
	ReadOnly *bool   `json:"readOnly,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (c *ConfigRequest) Validate() error {
	if c.Password != nil && len(*c.Password) > 2<<6 {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("password cannot be longer than 64 characters"))
	}
	return nil
}
