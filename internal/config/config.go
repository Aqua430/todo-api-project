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
	JWT
}

type HTTPServer struct {
	Port string `env:"PORT" env-default:"8080"`
}

type Database struct {
	URL string `env:"DATABASE_URL" env-required:"true"`
}

type JWT struct {
	Secret   string `env:"JWT_SECRET" env-required:"true"`
	TTLHours int    `env:"JWT_TTL_HOURS" env-default:"24"`
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
