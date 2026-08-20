package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `env:"ENV" env-default:"local"`
	HTTPServer
	Database
}

type HTTPServer struct {
	Port string `env:"PORT" env-default:"8080"`
}

type Database struct {
	URL string `env:"DATABASE_URL" env-required:"true"`
}

var (
	cfg  *Config
	once sync.Once
)

func MustLoad() *Config {
	once.Do(func() {
		cfg = &Config{}

		if err := cleanenv.ReadConfig(".env", cfg); err != nil {
			if err := cleanenv.ReadEnv(cfg); err != nil {
				log.Fatalf("Failed to read config: %v", err)
			}
		}
	})

	return cfg
}
