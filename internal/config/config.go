package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen struct {
		Addr   string
		BindIP string `env:"BIND_IP"`
		Port   string `env:"LISTEN_PORT"`
	}
	Postgresql struct {
		DSN      string
		Host     string `env:"PSQL_HOST"`
		Port     string `env:"PSQL_PORT"`
		Database string `env:"PSQL_NAME"`
		Username string `env:"PSQL_USER"`
		Password string `env:"PSQL_PASSWORD"`
	}
	Jwt struct {
		SecretKey string `env:"JWT_SECRET_KEY"`
	}
	CORS struct {
		AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
	}
}

var instance *Config
var once sync.Once

// GetConfig reads and parses the application's configuration file (.env).
// It uses double-check locking mechanism via sync.Once to ensure thread safety.
// The method initializes a Config object, populates it with values read from .env file,
// constructs the listen address and Postgres DSN, then logs success upon completion.
func GetConfig() *Config {
	once.Do(func() {
		log.Println("reading app configuration")
		instance = &Config{}
		err := cleanenv.ReadConfig("config.env", instance)
		if err != nil {
			log.Fatalln("read app configuration error")
		}
		instance.Listen.Addr = instance.Listen.BindIP + ":" + instance.Listen.Port
		instance.Postgresql.DSN = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
			instance.Postgresql.Username, instance.Postgresql.Password, instance.Postgresql.Host,
			instance.Postgresql.Port, instance.Postgresql.Database)
		log.Println("reading config OK")
	})
	return instance
}
