package infra

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Title                 string `yaml:"title"`
	Email                 string `yaml:"email"`
	PasswordToken         string `yaml:"passwordToken"`
	UrlPageChangePassword string `yaml:"urlPageChangePassword"`
	Server                struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		JWTToken string `yaml:"jwtToken"`
	}
	DataBase struct {
		Name string `yaml:"name"`
		URI  string `yaml:"uri"`
	}
	SendGrid struct {
		APIKey string `yaml:"apiKey"`
	}
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
