package daos

import (
	"context"
	"database/sql"

	"go-template/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// CreateRoleTx ...
func CreateRoleTx(role models.Role, ctx context.Context, tx *sql.Tx) (models.Role, error) {
	contextExecutor := GetContextExecutor(tx)

	err := role.Insert(ctx, contextExecutor, boil.Infer())
	return role, err
}

// CreateRoleTx ...
func CreateRole(role models.Role, ctx context.Context) (models.Role, error) {
	return CreateRoleTx(role, ctx, nil)
}

// FindRoleByID ...
func FindRoleByID(roleID int, ctx context.Context) (*models.Role, error) {
	contextExecutor := GetContextExecutor(nil)
	return models.FindRole(ctx, contextExecutor, roleID)
}
