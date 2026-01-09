package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen     ListenConfig
	Postgresql PostgresConfig
	Jwt        JwtConfig
	CORS       CORSConfig
}

type ListenConfig struct {
	BindIP string `env:"BIND_IP" env-default:"127.0.0.1"`
	Port   string `env:"LISTEN_PORT" env-default:"8080"`
}

func (l ListenConfig) Addr() string {
	return fmt.Sprintf("%s:%s", l.BindIP, l.Port)
}

type PostgresConfig struct {
	Host     string `env:"PSQL_HOST" env-required:"true"`
	Port     string `env:"PSQL_PORT" env-default:"5432"`
	Database string `env:"PSQL_NAME" env-required:"true"`
	Username string `env:"PSQL_USER" env-required:"true"`
	Password string `env:"PSQL_PASSWORD" env-required:"true"`
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		p.Username, p.Password, p.Host, p.Port, p.Database)
}

type JwtConfig struct {
	SecretKey string `env:"JWT_SECRET_KEY" env-required:"true"`
}

type CORSConfig struct {
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" env-separator:","`
}

func MustLoad(path string) *Config {
	var cfg Config

	log.Printf("Reading config from %s...", path)

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	log.Println("Config loaded successfully")
	return &cfg
}
