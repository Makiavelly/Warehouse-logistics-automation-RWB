package config

import "time"

type ServerConfig struct {
	Address            string        `yaml:"address"              env:"SERVER_ADDRESS"              env-default:":8080"`
	ReadTimeout        time.Duration `yaml:"read_timeout"         env:"SERVER_READ_TIMEOUT"         env-default:"15s"`
	WriteTimeout       time.Duration `yaml:"write_timeout"        env:"SERVER_WRITE_TIMEOUT"        env-default:"15s"`
	MaxShutdownTimeout time.Duration `yaml:"max_shutdown_timeout" env:"SERVER_MAX_SHUTDOWN_TIMEOUT" env-default:"10s"`
}