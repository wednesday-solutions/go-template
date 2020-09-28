package daos

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/wednesday-solutions/go-boiler/models"
)

// FindUserByUserName finds user by username
func FindUserByUserName(username string) (*models.User, error) {
	contextExecutor := getContextExecutor(nil)
	return models.Users(qm.Where(fmt.Sprintf("%s=?", models.UserColumns.Username), username)).One(context.Background(), contextExecutor)
}

// FindUserByEmail ...
func FindUserByEmail(email string) (*models.User, error) {
	contextExecutor := getContextExecutor(nil)
	return models.Users(qm.Where(fmt.Sprintf("%s=?", models.UserColumns.Email), email)).One(context.Background(), contextExecutor)
}

// FindUserByToken ...
func FindUserByToken(token string) (*models.User, error) {
	contextExecutor := getContextExecutor(nil)
	return models.Users(qm.Where(fmt.Sprintf("%s=?", models.UserColumns.Token), token)).One(context.Background(), contextExecutor)
}

// FindUserByID ...
func FindUserByID(userID int) (*models.User, error) {
	contextExecutor := getContextExecutor(nil)
	return models.FindUser(context.Background(), contextExecutor, userID)
}

// CreateUserTx ...
func CreateUserTx(user models.User, tx *sql.Tx) (models.User, error) {
	contextExecutor := getContextExecutor(tx)

	err := user.Insert(context.Background(), contextExecutor, boil.Infer())
	return user, err
}

// UpdateUserTx ...
func UpdateUserTx(user models.User, tx *sql.Tx) (models.User, error) {
	contextExecutor := getContextExecutor(tx)
	_, err := user.Update(context.Background(), contextExecutor, boil.Infer())
	return user, err
}

// DeleteUser ...
func DeleteUser(user models.User) (int64, error) {
	contextExecutor := getContextExecutor(nil)
	rowsAffected, err := user.Delete(context.Background(), contextExecutor)
	return rowsAffected, err
}
