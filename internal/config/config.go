package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config stores whole application settings.
	Config struct {
		// DBMS configures which data source type application will using.
		DBMS string `yaml:"dbms" env:"DBMS" env-default:"sqlite"`

		Logger Logger `yaml:"logger"`
		SQLite SQLite `yaml:"sqlite"`
		HTTP   HTTP   `yaml:"http"`
	}

	// Logger stores log settings.
	Logger struct {
		// Level is log level. Possible values: debug, info, warn, error
		Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`
	}

	// SQLite stores sqlite3 database settings.
	SQLite struct {
		// DSN is a destination string to database file.
		DSN string `yaml:"dsn" env:"SQLITE_DSN" env-default:"homeset.sqlite"`
	}

	// HTTP stores http server settings.
	HTTP struct {
		JWTToken string `yaml:"jwtToken" env:"JWT_TOKEN" env-required:"true"`
		// Port is a http server port without :.
		Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
		// ReadTimeout is a server read timeout.
		ReadTimeout time.Duration `yaml:"readTimeout" env:"HTTP_READ_TIMEOUT" env-default:"10s"`
		// WriteTimeout is a server write timeout.
		WriteTimeout time.Duration `yaml:"writeTimeout" env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
		// ShutdownTimeout is a server shutdown timeout.
		ShutdownTimeout time.Duration `yaml:"shutdownTimeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"10s"`
	}
)

// New creates a new Config.
func New(configPath string) (Config, error) {
	cfg := Config{}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadConfig: %w", err)
	}

	return cfg, nil
}
