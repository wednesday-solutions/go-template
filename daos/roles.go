package daos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/wednesday-solutions/go-template/models"
)

// CreateRoleTx ...
func CreateRoleTx(role models.Role, tx *sql.Tx) (models.Role, error) {
	contextExecutor := getContextExecutor(tx)
	boil.DebugMode = true
	err := role.Insert(context.Background(), contextExecutor, boil.Infer())
	fmt.Print(err)
	return role, err
}

// FindRoleByID ...
func FindRoleByID(roleID int) (*models.Role, error) {
	contextExecutor := getContextExecutor(nil)
	return models.FindRole(context.Background(), contextExecutor, roleID)
}
