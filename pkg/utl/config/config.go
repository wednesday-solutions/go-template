package config

import (
	"github.com/wednesday-solutions/go-template/pkg/utl/convert"
	"os"
)

// Load returns Configuration struct
func Load() (*Configuration, error) {
	cfg := &Configuration{
		Server: &Server{
			Port:         os.Getenv("SERVER_PORT"),
			Debug:        convert.StringToBool(os.Getenv("SERVER_DEBUG")),
			ReadTimeout:  convert.StringToInt(os.Getenv("SERVER_READ_TIMEOUT")),
			WriteTimeout: convert.StringToInt(os.Getenv("SERVER_WRITE_TIMEOUT")),
		},
		DB: &Database{
			LogQueries: convert.StringToBool(os.Getenv("DB_LOG_QUERIES")),
			Timeout:    convert.StringToInt(os.Getenv("DB_TIMEOUT_SECONDS")),
		},
		JWT: &JWT{
			MinSecretLength:  convert.StringToInt(os.Getenv("JWT_MIN_SECRET_LENGTH")),
			DurationMinutes:  convert.StringToInt(os.Getenv("JWT_DURATION_MINUTES")),
			RefreshDuration:  convert.StringToInt(os.Getenv("JWT_REFRESH_DURATION")),
			MaxRefresh:       convert.StringToInt(os.Getenv("JWT_MAX_REFRESH")),
			SigningAlgorithm: os.Getenv("JWT_SIGNING_ALGORITHM"),
		},
		App: &Application{
			MinPasswordStr: convert.StringToInt(os.Getenv("APP_MIN_PASSWORD_STR")),
		},
	}
	return cfg, nil
}

// Configuration holds data necessary for configuring application
type Configuration struct {
	Server *Server      `yaml:"server,omitempty"`
	DB     *Database    `yaml:"database,omitempty"`
	JWT    *JWT         `yaml:"jwt,omitempty"`
	App    *Application `yaml:"application,omitempty"`
}

// Database holds data necessary for database configuration
type Database struct {
	LogQueries bool `yaml:"log_queries,omitempty"`
	Timeout    int  `yaml:"timeout_seconds,omitempty"`
}

// Server holds data necessary for server configuration
type Server struct {
	Port         string `yaml:"port,omitempty"`
	Debug        bool   `yaml:"debug,omitempty"`
	ReadTimeout  int    `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeout int    `yaml:"write_timeout_seconds,omitempty"`
}

// JWT holds data necessary for JWT configuration
type JWT struct {
	MinSecretLength  int    `yaml:"min_secret_length,omitempty"`
	DurationMinutes  int    `yaml:"duration_minutes,omitempty"`
	RefreshDuration  int    `yaml:"refresh_duration_minutes,omitempty"`
	MaxRefresh       int    `yaml:"max_refresh_minutes,omitempty"`
	SigningAlgorithm string `yaml:"signing_algorithm,omitempty"`
}

// Application holds application configuration details
type Application struct {
	MinPasswordStr int `yaml:"min_password_strength,omitempty"`
}
