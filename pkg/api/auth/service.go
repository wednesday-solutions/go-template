package auth

import (
	"database/sql"
	"github.com/labstack/echo"
	"github.com/wednesday-solutions/go-boiler/models"

	"github.com/wednesday-solutions/go-boiler"
)

// New creates new iam service
func New(db *sql.DB, j TokenGenerator, sec Securer) Auth {
	return Auth{
		db:  db,
		tg:  j,
		sec: sec,
	}
}

// Initialize initializes auth application service
func Initialize(db *sql.DB, j TokenGenerator, sec Securer) Auth {
	return New(db, j, sec)
}

// Service represents auth service interface
type Service interface {
	Authenticate(echo.Context, string, string) (goboiler.AuthToken, error)
	Refresh(echo.Context, string) (string, error)
	Me(echo.Context) (*models.User, error)
}

// Auth represents auth application service
type Auth struct {
	db  *sql.DB
	tg  TokenGenerator
	sec Securer
}

// TokenGenerator represents token generator (jwt) interface
type TokenGenerator interface {
	GenerateToken(*models.User) (string, error)
}

// Securer represents security interface
type Securer interface {
	HashMatchesPassword(string, string) bool
	Token(string) string
}
