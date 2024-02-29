package daos

import (
	"context"
	"database/sql"
	"go-template/models"
	"strconv"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func FindPostWithId(postId string, ctx context.Context) (*models.Post, error) {
	contextExecutor := GetContextExecutor(nil)
	postIdInt, _ := strconv.Atoi(postId)
	return models.FindPost(ctx, contextExecutor, postIdInt)
}
func FindAllPostsWithCount(queryMods []qm.QueryMod, ctx context.Context) (models.PostSlice, int64, error) {
	contextExecutor := GetContextExecutor(nil)
	// Get all posts that match the given query modifiers.
	posts, err := models.Posts(queryMods...).All(ctx, contextExecutor)
	if err != nil {
		return models.PostSlice{}, 0, err
	}
	// Get the count of all users that match the given query modifiers.
	queryMods = append(queryMods, qm.Offset(0))
	count, err := models.Posts(queryMods...).Count(ctx, contextExecutor)
	return posts, count, err
}

// CreatePostTx creates a new post in the database.
func CreatePostTx(post models.Post, ctx context.Context, tx *sql.Tx) (models.Post, error) {
	contextExecutor := GetContextExecutor(tx)
	err := post.Insert(ctx, contextExecutor, boil.Infer())
	return post, err
}

// CreatePost creates a new post in the database.
func CreatePost(post models.Post, ctx context.Context) (models.Post, error) {
	return CreatePostTx(post, ctx, nil)
}

// UpdatePostTx updates an post in the database.
func UpdatePostTx(post models.Post, ctx context.Context, tx *sql.Tx) (models.Post, error) {
	contextExecutor := GetContextExecutor(tx)
	// Update the post in the database.
	_, err := post.Update(ctx, contextExecutor, boil.Infer())
	return post, err
}

// UpdatePost updates an post in the database.
func UpdatePost(post models.Post, ctx context.Context) (models.Post, error) {
	return UpdatePostTx(post, ctx, nil)
}

func DeletePost(post models.Post, ctx context.Context) (int64, error) {
	contextExecutor := GetContextExecutor(nil)
	// Delete the post from the database.
	rowsAffected, err := post.Delete(ctx, contextExecutor)
	return rowsAffected, err
}
