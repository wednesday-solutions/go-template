package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"

	"github.com/wednesday-solution/go-boiler"
	"github.com/wednesday-solution/go-boiler/pkg/api/user/platform/pgsql"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, goboiler.User) (goboiler.User, error)
	List(echo.Context, goboiler.Pagination) ([]goboiler.User, error)
	View(echo.Context, int) (goboiler.User, error)
	Delete(echo.Context, int) error
	Update(echo.Context, Update) (goboiler.User, error)
}

// New creates new user application service
func New(db *pg.DB, udb UDB, rbac RBAC, sec Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *User {
	return New(db, pgsql.User{}, rbac, sec)
}

// User represents user application service
type User struct {
	db   *pg.DB
	udb  UDB
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// UDB represents user repository interface
type UDB interface {
	Create(orm.DB, goboiler.User) (goboiler.User, error)
	View(orm.DB, int) (goboiler.User, error)
	List(orm.DB, *goboiler.ListQuery, goboiler.Pagination) ([]goboiler.User, error)
	Update(orm.DB, goboiler.User) error
	Delete(orm.DB, goboiler.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) goboiler.AuthUser
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, goboiler.AccessRole, int, int) error
	IsLowerRole(echo.Context, goboiler.AccessRole) error
}
