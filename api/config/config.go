package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"app"`

	Server struct {
		BaseURL        string `yaml:"base_url"`
		Port           uint16 `yaml:"port"`
		AllowedOrigins string `yaml:"origins"`
		Metrics        struct {
			Enabled  bool   `yaml:"enabled"`
			Route    string `yaml:"route"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		} `yaml:"metrics"`
	} `yaml:"server"`

	Cache struct {
		Host string `yaml:"host"`
		Port uint16 `yaml:"port"`
	} `yaml:"cache"`

	Session `yaml:"session"`
	Github  `yaml:"github"`
}

type Session struct {
	DefaultUsers int    `yaml:"users"`
	MaxUsers     int    `yaml:"max_users"`
	RPM          uint16 `yaml:"rpm"`
}

type Github struct {
	Enabled           bool     `yaml:"enabled"`
	ClientId          string   `yaml:"client_id"`
	ClientSecret      string   `yaml:"client_secret"`
	RedirectURI       string   `yaml:"redirect_uri"`
	Scope             []string `yaml:"scope"`
	Emails            []string `yaml:"whitelisted_emails"`
	WhitelistedEmails map[string]struct{}
}

func (gh *Github) parseEmails() {
	whitelisted := make(map[string]struct{}, len(gh.Emails))
	for _, e := range gh.Emails {
		whitelisted[e] = struct{}{}
	}
	gh.WhitelistedEmails = whitelisted
}

func New(path string) (*Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Configuration{}
	err = yaml.NewDecoder(file).Decode(cfg)
	cfg.Github.parseEmails()
	return cfg, err
}
