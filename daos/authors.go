package daos

import (
	"context"
	"database/sql"
	"go-template/models"
	"strconv"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func FindAuthorWithId(authorId string, ctx context.Context) (*models.Author, error) {
	contextExecutor := GetContextExecutor(nil)
	authorIdInt, _ := strconv.Atoi(authorId)
	return models.FindAuthor(ctx, contextExecutor, authorIdInt)
}

// FindAllAuthorsWithCount returns all authors and their count, filtered by the given query modifiers.
func FindAllAuthorsWithCount(queryMods []qm.QueryMod, ctx context.Context) (models.AuthorSlice, int64, error) {
	contextExecutor := GetContextExecutor(nil)
	// Get all authors that match the given query modifiers.
	authors, err := models.Authors(queryMods...).All(ctx, contextExecutor)
	if err != nil {
		return models.AuthorSlice{}, 0, err
	}
	// Get the count of all users that match the given query modifiers.
	queryMods = append(queryMods, qm.Offset(0))
	count, err := models.Authors(queryMods...).Count(ctx, contextExecutor)
	return authors, count, err
}

// CreateAuthorTx creates a new author in the database, using a transaction.
func CreateAuthorTx(author models.Author, ctx context.Context, tx *sql.Tx) (models.Author, error) {
	contextExecutor := GetContextExecutor(tx)
	err := author.Insert(ctx, contextExecutor, boil.Infer())
	return author, err
}

// CreateAuthor creates a new author in the database.
func CreateAuthor(author models.Author, ctx context.Context) (models.Author, error) {
	return CreateAuthorTx(author, ctx, nil)
}

// UpdateAuthorTx updates an author in the database, using a transaction.
func UpdateAuthorTx(author models.Author, ctx context.Context, tx *sql.Tx) (models.Author, error) {
	contextExecutor := GetContextExecutor(tx)
	// Update the author in the database.
	_, err := author.Update(ctx, contextExecutor, boil.Infer())
	return author, err
}

// UpdateAuthor updates an author in the database.
func UpdateAuthor(author models.Author, ctx context.Context) (models.Author, error) {
	return UpdateAuthorTx(author, ctx, nil)
}

// DeleteAuthor deletes an author from the database.
func DeleteAuthor(author models.Author, ctx context.Context) (int64, error) {
	contextExecutor := GetContextExecutor(nil)
	// Delete the author from the database.
	rowsAffected, err := author.Delete(ctx, contextExecutor)
	return rowsAffected, err
}
