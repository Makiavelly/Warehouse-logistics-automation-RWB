package config

type MLConfig struct {
	Address  string `yaml:"address"   env:"ML_ADDRESS"    env-default:"ml-service:50051"`
	UseGRPC  bool   `yaml:"use_grpc"  env:"ML_USE_GRPC"   env-default:"true"`
}