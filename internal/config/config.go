package config

import (
	"fmt"
	"os"

	"go-template/pkg/utl/convert"
)

// Load returns Configuration struct
func Load() (*Configuration, error) {
	cfg := &Configuration{
		Server: &Server{
			Port:         fmt.Sprintf(":%d", convert.StringToInt(os.Getenv("SERVER_PORT"))),
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
	if len(os.Getenv("SERVER_PORT")) == 0 {
		return nil, fmt.Errorf("error loading port from .env ")
	}
	if len(os.Getenv("DB_TIMEOUT_SECONDS")) == 0 {
		return nil, fmt.Errorf("error loading db timeout from .env ")
	}
	if len(os.Getenv("JWT_MIN_SECRET_LENGTH")) == 0 {
		return nil, fmt.Errorf("error loading jwt min secret length from .env ")
	}
	if len(os.Getenv("APP_MIN_PASSWORD_STR")) == 0 {
		return nil, fmt.Errorf("error loading application password string from .env ")
	}
	if len(os.Getenv("SERVER_READ_TIMEOUT")) == 0 || len(os.Getenv("SERVER_WRITE_TIMEOUT")) == 0 {
		return nil, fmt.Errorf("error loading server timeout from .env ")
	}
	return cfg, nil
}

// Configuration holds data necessary for configuring application
type Configuration struct {
	Server *Server      `json:"server,omitempty"`
	DB     *Database    `json:"database,omitempty"`
	JWT    *JWT         `json:"jwt,omitempty"`
	App    *Application `json:"application,omitempty"`
}

// Database holds data necessary for database configuration
type Database struct {
	LogQueries bool `json:"log_queries,omitempty"`
	Timeout    int  `json:"timeout_seconds,omitempty"`
}

// Server holds data necessary for server configuration
type Server struct {
	Port         string `json:"port"                  validate:"required"`
	Debug        bool   `json:"debug"                 validate:"required"`
	ReadTimeout  int    `json:"read_timeout_seconds"  validate:"required"`
	WriteTimeout int    `json:"write_timeout_seconds" validate:"required"`
}

// JWT holds data necessary for JWT configuration
type JWT struct {
	MinSecretLength  int    `json:"min_secret_length"                  validate:"required"`
	DurationMinutes  int    `json:"duration_minutes,omitempty"`
	RefreshDuration  int    `json:"refresh_duration_minutes,omitempty"`
	MaxRefresh       int    `json:"max_refresh_minutes,omitempty"`
	SigningAlgorithm string `json:"signing_algorithm"                  validate:"required"`
}

// Application holds application configuration details
type Application struct {
	MinPasswordStr int `json:"min_password_strength" validate:"required"`
}
