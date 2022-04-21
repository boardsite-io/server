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

	Cache struct {
		Host string `yaml:"host"`
		Port uint16 `yaml:"port"`
	} `yaml:"cache"`

	Server  `yaml:"server"`
	Session `yaml:"session"`
}

type Server struct {
	BaseURL        string `yaml:"base_url"`
	Port           uint16 `yaml:"port"`
	AllowedOrigins string `yaml:"origins"`
	RPM            uint16 `yaml:"rpm"`
	Metrics        struct {
		Enabled  bool   `yaml:"enabled"`
		Route    string `yaml:"route"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"metrics"`
}

type Session struct {
	MaxUsers int  `yaml:"max_users" json:"maxUsers"`
	ReadOnly bool `yaml:"read_only" json:"readOnly"`
}

func New(path string) (*Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Configuration{}
	err = yaml.NewDecoder(file).Decode(cfg)
	return cfg, err
}
