package config

type DBConfig struct {
	Address string `yaml:"address" env:"DB_ADDRESS" env-default:"postgres://postgres:password@localhost:5432/logistics?sslmode=disable"`
}