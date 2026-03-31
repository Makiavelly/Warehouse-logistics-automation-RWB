package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	App    AppConfig    `yaml:"app"`
	Auth   AuthConfig   `yaml:"auth"`
	DB     DBConfig     `yaml:"db"`
	Server ServerConfig `yaml:"server"`
	ML     MLConfig     `yaml:"ml"`
}

func MustLoadConfig(configPath string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}
	return &cfg
}