package password

import (
	"database/sql"
	"github.com/labstack/echo"
)

// Service represents password application interface
type Service interface {
	Change(echo.Context, int, string, string) error
}

// New creates new password application service
func New(db   *sql.DB, sec Securer) Password {
	return Password{
		db:   db,
		sec:  sec,
	}
}

// Initialize initialises password application service with defaults
func Initialize(db   *sql.DB, sec Securer) Password {
	return New(db, sec)
}

// Password represents password application service
type Password struct {
	db   *sql.DB
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
	HashMatchesPassword(string, string) bool
	Password(string, ...string) bool
}