package config

import (
	"github.com/aliykh/log"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/qiangxue/go-env"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	defaultServerPort = 7077
)

// Config represents an application configuration.
type Config struct {
	// the server port. Defaults to 7077
	ServerPort int `yaml:"server_port" env:"SERVER_PORT"`

	// mongo connection uri string
	MongoAddr string `yaml:"mongo_addr" env:"MONGO_ADDR,secret"`

	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`

	DatabaseName string `yaml:"db_name" env:"DATABASE_NAME"`
}

// Validate validates the application configuration.
func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.MongoAddr, validation.Required),
		validation.Field(&c.DatabaseName, validation.Required),
	)
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(file string, logger log.Logger) (*Config, error) {
	// default config
	c := Config{
		ServerPort: defaultServerPort,
	}

	// load from YAML config file
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	// load from environment variables prefixed with "APP_"
	if err = env.New("APP_", nil).Load(&c); err != nil {
		return nil, err
	}

	// validation
	if err = c.Validate(); err != nil {
		return nil, err
	}

	return &c, err
}
