package config

import "time"

type AuthConfig struct {
	SecretKey  string        `yaml:"secret_key"  env:"AUTH_SECRET_KEY"  env-default:"change-me-in-production"`
	TTLHours   time.Duration `yaml:"ttl_hours"   env:"AUTH_TTL_HOURS"   env-default:"72"`
	APIKey     string        `yaml:"api_key"     env:"AUTH_API_KEY"     env-default:"internal-api-key"`
}