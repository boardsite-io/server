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
		Host           string `yaml:"host"`
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

	Session struct {
		MaxUsers int    `yaml:"users"`
		RPM      uint16 `yaml:"rpm"`
	} `yaml:"session"`
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
