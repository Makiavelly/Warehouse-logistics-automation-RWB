package config

type AppConfig struct {
	LogLevel string `yaml:"log_level" env:"APP_LOG_LEVEL" env-default:"info"`
}