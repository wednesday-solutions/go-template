package user

import (
	"database/sql"
	"github.com/labstack/echo"
	"github.com/wednesday-solution/go-boiler"
	"github.com/wednesday-solution/go-boiler/models"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, models.User) (models.User, error)
	List(echo.Context, goboiler.Pagination) (models.UserSlice, error)
	View(echo.Context, int) (*models.User, error)
	Delete(echo.Context, int) error
	Update(echo.Context, models.User) (models.User, error)
}

// New creates new user application service
func New(db *sql.DB, sec Securer) *User {
	return &User{db: db, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(db *sql.DB, sec Securer) *User {
	return New(db, sec)
}

// User represents user application service
type User struct {
	db  *sql.DB
	sec Securer
}

// Securer represents securityInitialize interface
type Securer interface {
	Hash(string) string
}
