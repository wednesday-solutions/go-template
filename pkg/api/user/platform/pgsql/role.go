package pgsql

import (
	"github.com/go-pg/pg/orm"
	goboiler "github.com/wednesday-solution/go-boiler"
)

type Role struct{}

func (u Role) View(db orm.DB, id int) (goboiler.Role, error) {
	var role goboiler.Role
	sql := `SELECT * from roles WHERE ("role""."id" = ?)`
	_, err := db.QueryOne(&role, sql, id)
	return role, err
}