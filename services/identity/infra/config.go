package infra

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Title  string `yaml:"title"`
	Email  string `yaml:"email"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	DataBase struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	SendGrid struct {
		APIKey string `yaml:"apiKey"`
	}
	APIKey string `yaml:"key"`
}

func NewConfig(filename string) (*Config, error) {
	config := &Config{}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, config)

	return config, err
}

func (cfg Config) ServerAddress() string {
	return cfg.Server.Host + ":" + cfg.Server.Port
}

func (cfg Config) DataBaseAddress() string {
	return cfg.DataBase.Host + ":" + cfg.DataBase.Port
}
