package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type DBConfig struct {
	Connection string `yaml:"connection" env:"DB_CONNECTION" env-required:"true"`
	User       string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password   string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Database   string `yaml:"database" env:"DB_DATABASE" env-required:"true"`
	Port       int    `yaml:"port" env:"DB_PORT" env-required:"true"`
}

type GRPCConfig struct {
	Port int `yaml:"port" env:"GRPC_PORT" env-required:"true"`
}

type Config struct {
	Env    string     `yaml:"env" env:"ENV" env-default:"local"`
	DB     DBConfig   `yaml:"db"`
	GRPC   GRPCConfig `yaml:"grpc"`
	Secret string     `yaml:"secret" env:"SECRET" env-required:"true"`
}

func MustLoad() *Config {
	configPath := resolveConfigPath()
	if configPath == "" {
		panic("unable to resolve config path")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to parse config: " + err.Error())
	}

	return &cfg
}

func resolveConfigPath() string {
	var path string

	flag.StringVar(&path, "config path", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
