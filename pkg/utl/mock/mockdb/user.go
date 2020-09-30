package mockdb

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/wednesday-solutions/go-template/models"
)

// User database mock
type User struct {
	CreateFn         func(orm.DB, models.User) (models.User, error)
	ViewFn           func(orm.DB, int) (models.User, error)
	FindByUsernameFn func(orm.DB, string) (models.User, error)
	FindByTokenFn    func(orm.DB, string) (models.User, error)
	DeleteFn         func(orm.DB, models.User) error
	UpdateFn         func(orm.DB, models.User) error
}

// Create mock
func (u *User) Create(db orm.DB, usr models.User) (models.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *User) View(db orm.DB, id int) (models.User, error) {
	return u.ViewFn(db, id)
}

// FindByUsername mock
func (u *User) FindByUsername(db orm.DB, uname string) (models.User, error) {
	return u.FindByUsernameFn(db, uname)
}

// FindByToken mock
func (u *User) FindByToken(db orm.DB, token string) (models.User, error) {
	return u.FindByTokenFn(db, token)
}

// Delete mock
func (u *User) Delete(db orm.DB, usr models.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *User) Update(db orm.DB, usr models.User) error {
	return u.UpdateFn(db, usr)
}
