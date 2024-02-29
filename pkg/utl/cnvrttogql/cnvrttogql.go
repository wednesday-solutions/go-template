package cnvrttogql

import (
	"context"
	graphql "go-template/gqlmodels"
	"go-template/internal/constants"
	"go-template/models"
	"go-template/pkg/utl/convert"
	"strconv"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// UsersToGraphQlUsers converts array of type models.User into array of pointer type graphql.User
func UsersToGraphQlUsers(u models.UserSlice, count int) []*graphql.User {
	var r []*graphql.User
	for _, e := range u {
		r = append(r, UserToGraphQlUser(e, count))
	}
	return r
}

// UserToGraphQlUser converts type models.User into pointer type graphql.User
func UserToGraphQlUser(u *models.User, count int) *graphql.User {
	count++
	if u == nil {
		return nil
	}
	var role *models.Role
	if count <= constants.MaxDepth {
		u.L.LoadRole(context.Background(), boil.GetContextDB(), true, u, nil) //nolint:errcheck
		if u.R != nil {
			role = u.R.Role
		}
	}

	return &graphql.User{
		ID:        strconv.Itoa(u.ID),
		FirstName: convert.NullDotStringToPointerString(u.FirstName),
		LastName:  convert.NullDotStringToPointerString(u.LastName),
		Username:  convert.NullDotStringToPointerString(u.Username),
		Email:     convert.NullDotStringToPointerString(u.Email),
		Mobile:    convert.NullDotStringToPointerString(u.Mobile),
		Address:   convert.NullDotStringToPointerString(u.Address),
		Active:    convert.NullDotBoolToPointerBool(u.Active),
		Role:      RoleToGraphqlRole(role, count),
	}
}

func RoleToGraphqlRole(r *models.Role, count int) *graphql.Role {
	count++
	if r == nil {
		return nil
	}
	var users models.UserSlice
	if count <= constants.MaxDepth {
		r.L.LoadUsers(context.Background(), boil.GetContextDB(), true, r, nil) //nolint:errcheck
		if r.R != nil {
			users = r.R.Users
		}
	}

	return &graphql.Role{
		ID:          strconv.Itoa(r.ID),
		AccessLevel: r.AccessLevel,
		Name:        r.Name,
		UpdatedAt:   convert.NullDotTimeToPointerInt(r.UpdatedAt),
		CreatedAt:   convert.NullDotTimeToPointerInt(r.CreatedAt),
		DeletedAt:   convert.NullDotTimeToPointerInt(r.DeletedAt),
		Users:       UsersToGraphQlUsers(users, count),
	}
}

func AuthorsToGraphQlAuthors(u models.AuthorSlice) []*graphql.Author {
	var r []*graphql.Author
	for _, e := range u {
		r = append(r, AuthorToGraphQlAuthor(e))
	}
	return r
}
func AuthorToGraphQlAuthor(u *models.Author) *graphql.Author {
	if u == nil {
		return nil
	}
	return &graphql.Author{
		ID:        strconv.Itoa(u.ID),
		FirstName: convert.NullDotStringToPointerString(u.FirstName),
		LastName:  convert.NullDotStringToPointerString(u.LastName),
		Email:     convert.NullDotStringToPointerString(u.Email),
	}
}

// PostsToGraphQlPosts converts type models.Post into pointer type graphql.Post
func PostsToGraphQlPosts(u models.PostSlice, count int64) []*graphql.Post {
	var r []*graphql.Post
	for _, e := range u {
		r = append(r, PostToGraphQlPost(e, count))
	}
	return r
}

// PostToGraphQlPost converts type models.Post into pointer type graphql.Post
func PostToGraphQlPost(u *models.Post, count int64) *graphql.Post {
	count++
	if u == nil {
		return nil
	}
	var author *models.Author
	if count <= constants.MaxDepth {
		u.L.LoadAuthor(context.Background(), boil.GetContextDB(), true, u, nil) //nolint:errcheck
		if u.R != nil {
			author = u.R.Author
		}
	}

	return &graphql.Post{
		ID:     strconv.Itoa(u.ID),
		Author: AuthorToGraphQlAuthor(author),
		Post:   &u.Post,
	}
}
