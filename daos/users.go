package daos

import (
	"context"
	"database/sql"
	"fmt"

	"go-template/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// FindUserByUserName finds user by username
func FindUserByUserName(username string, ctx context.Context) (*models.User, error) {
	contextExecutor := GetContextExecutor(nil)
	return models.Users(qm.Where(fmt.Sprintf("%s=?", models.UserColumns.Username), username)).
		One(ctx, contextExecutor)
}

// FindUserByEmail ...
func FindUserByEmail(email string, ctx context.Context) (*models.User, error) {
	contextExecutor := GetContextExecutor(nil)
	return models.Users(qm.Where(fmt.Sprintf("%s=?", models.UserColumns.Email), email)).
		One(ctx, contextExecutor)
}

// FindUserByToken ...
func FindUserByToken(token string, ctx context.Context) (*models.User, error) {
	contextExecutor := GetContextExecutor(nil)
	return models.Users(qm.Where(fmt.Sprintf("%s=?", models.UserColumns.Token), token)).
		One(ctx, contextExecutor)
}

// FindUserByID ...
func FindUserByID(userID int, ctx context.Context) (*models.User, error) {
	contextExecutor := GetContextExecutor(nil)
	return models.FindUser(ctx, contextExecutor, userID)
}

// CreateUserTx ...
func CreateUserTx(user models.User, ctx context.Context, tx *sql.Tx) (models.User, error) {
	contextExecutor := GetContextExecutor(tx)

	err := user.Insert(ctx, contextExecutor, boil.Infer())
	return user, err
}

// CreateUser
func CreateUser(user models.User, ctx context.Context) (models.User, error) {
	return CreateUserTx(user, ctx, nil)
}

// UpdateUserTx ...
func UpdateUserTx(user models.User, ctx context.Context, tx *sql.Tx) (models.User, error) {
	contextExecutor := GetContextExecutor(tx)
	_, err := user.Update(ctx, contextExecutor, boil.Infer())
	return user, err
}

// UpdateUserTx ...
func UpdateUser(user models.User, ctx context.Context) (models.User, error) {
	return UpdateUserTx(user, ctx, nil)
}

// DeleteUser ...
func DeleteUser(user models.User, ctx context.Context) (int64, error) {
	contextExecutor := GetContextExecutor(nil)
	rowsAffected, err := user.Delete(ctx, contextExecutor)
	return rowsAffected, err
}

// FindAllUsersWithCount ... This will get all the users that match the queryMod filter and also return the count
func FindAllUsersWithCount(queryMods []qm.QueryMod, ctx context.Context) (models.UserSlice, int64, error) {
	contextExecutor := GetContextExecutor(nil)
	users, err := models.Users(queryMods...).All(ctx, contextExecutor)
	if err != nil {
		return models.UserSlice{}, 0, err
	}
	queryMods = append(queryMods, qm.Offset(0))
	count, err := models.Users(queryMods...).Count(ctx, contextExecutor)
	return users, count, err

}
