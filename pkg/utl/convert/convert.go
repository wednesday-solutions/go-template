package convert

import (
	"context"
	"strconv"

	gotemplate "go-template"
	graphql "go-template/gqlmodels"
	"go-template/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// StringToPointerString returns pointer string value
func StringToPointerString(v string) *string {
	return &v
}

// StringToInt converts string to integer
func StringToInt(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}

// StringToBool converts string to boolean
func StringToBool(v string) bool {
	i, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return i
}

// NullDotStringToPointerString converts nullable string to its pointer value
func NullDotStringToPointerString(v null.String) *string {
	return v.Ptr()
}

// NullDotStringToString converts nullable string to its value
func NullDotStringToString(v null.String) string {
	if v.Ptr() == nil {
		return ""
	}
	return *v.Ptr()
}

// NullDotIntToInt converts nullable int to its value
func NullDotIntToInt(v null.Int) int {
	if v.Ptr() == nil {
		return 0
	}
	return *v.Ptr()
}

// NullDotBoolToPointerBool converts nullable boolean to its pointer value
func NullDotBoolToPointerBool(v null.Bool) *bool {
	return v.Ptr()
}

// PointerStringToNullDotInt converts pointer string to nullable integer if present else returns default nullable value
func PointerStringToNullDotInt(s *string) null.Int {
	if s == nil {
		return null.IntFrom(0)
	}
	v := *s
	i, err := strconv.Atoi(v)
	if err != nil {
		return null.IntFrom(0)
	}
	return null.IntFrom(i)
}

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
	if count <= gotemplate.MaxDepth {
		u.L.LoadRole(context.Background(), boil.GetContextDB(), true, u, nil) //nolint:errcheck
		if u.R != nil {
			role = u.R.Role
		}
	}

	return &graphql.User{
		ID:        strconv.Itoa(u.ID),
		FirstName: NullDotStringToPointerString(u.FirstName),
		LastName:  NullDotStringToPointerString(u.LastName),
		Username:  NullDotStringToPointerString(u.Username),
		Email:     NullDotStringToPointerString(u.Email),
		Mobile:    NullDotStringToPointerString(u.Mobile),
		Address:   NullDotStringToPointerString(u.Address),
		Active:    NullDotBoolToPointerBool(u.Active),
		Role:      RoleToGraphqlRole(role, count),
	}
}

func RoleToGraphqlRole(r *models.Role, count int) *graphql.Role {
	count++
	if r == nil {
		return nil
	}
	var users models.UserSlice
	if count <= gotemplate.MaxDepth {
		r.L.LoadUsers(context.Background(), boil.GetContextDB(), true, r, nil) //nolint:errcheck
		if r.R != nil {
			users = r.R.Users
		}
	}

	return &graphql.Role{
		ID:          strconv.Itoa(r.ID),
		AccessLevel: r.AccessLevel,
		Name:        r.Name,
		UpdatedAt:   NullDotTimeToPointerInt(r.UpdatedAt),
		CreatedAt:   NullDotTimeToPointerInt(r.CreatedAt),
		DeletedAt:   NullDotTimeToPointerInt(r.DeletedAt),
		Users:       UsersToGraphQlUsers(users, count),
	}
}
func NullDotTimeToPointerInt(t null.Time) *int {
	var i int
	if t.Valid {
		i = int(t.Time.UnixMilli())
		return &i
	}
	return nil
}
