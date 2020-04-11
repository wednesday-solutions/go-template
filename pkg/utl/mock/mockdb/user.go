package mockdb

import (
	"github.com/go-pg/pg/v9/orm"

	"github.com/wednesday-solutions/go-boiler"
)

// User database mock
type User struct {
	CreateFn         func(orm.DB, goboiler.User) (goboiler.User, error)
	ViewFn           func(orm.DB, int) (goboiler.User, error)
	FindByUsernameFn func(orm.DB, string) (goboiler.User, error)
	FindByTokenFn    func(orm.DB, string) (goboiler.User, error)
	ListFn           func(orm.DB, *goboiler.ListQuery, goboiler.Pagination) ([]goboiler.User, error)
	DeleteFn         func(orm.DB, goboiler.User) error
	UpdateFn         func(orm.DB, goboiler.User) error
}

// Create mock
func (u *User) Create(db orm.DB, usr goboiler.User) (goboiler.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *User) View(db orm.DB, id int) (goboiler.User, error) {
	return u.ViewFn(db, id)
}

// FindByUsername mock
func (u *User) FindByUsername(db orm.DB, uname string) (goboiler.User, error) {
	return u.FindByUsernameFn(db, uname)
}

// FindByToken mock
func (u *User) FindByToken(db orm.DB, token string) (goboiler.User, error) {
	return u.FindByTokenFn(db, token)
}

// List mock
func (u *User) List(db orm.DB, lq *goboiler.ListQuery, p goboiler.Pagination) ([]goboiler.User, error) {
	return u.ListFn(db, lq, p)
}

// Delete mock
func (u *User) Delete(db orm.DB, usr goboiler.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *User) Update(db orm.DB, usr goboiler.User) error {
	return u.UpdateFn(db, usr)
}
