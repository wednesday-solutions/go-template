package convert

import (
	"github.com/volatiletech/null"
	graphql "github.com/wednesday-solutions/go-boiler/graphql_models"
	"github.com/wednesday-solutions/go-boiler/models"
	"strconv"
)

// StringToPointerString ...
func StringToPointerString(v string) *string {
	return &v
}

// StringToInt ...
func StringToInt(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}

// StringToBool ...
func StringToBool(v string) bool {
	i, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return i
}

// NullDotStringToPointerString ...
func NullDotStringToPointerString(v null.String) *string {
	return v.Ptr()
}

// NullDotBoolToPointerBool ...
func NullDotBoolToPointerBool(v null.Bool) *bool {
	return v.Ptr()
}

// NullDotIntToPointerInt ...
func NullDotIntToPointerInt(v null.Int) *int {
	return v.Ptr()
}

// PointerStringToNullDotInt ...
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

// UsersToGraphQlUsers ...
func UsersToGraphQlUsers(u models.UserSlice) []*graphql.User {
	var r []*graphql.User
	for _, e := range u {
		r = append(r, UserToGraphQlUser(e))
	}
	return r
}

// UserToGraphQlUser ...
func UserToGraphQlUser(u *models.User) *graphql.User {
	return &graphql.User{
		ID:        strconv.Itoa(u.ID),
		FirstName: NullDotStringToPointerString(u.FirstName),
		LastName:  NullDotStringToPointerString(u.FirstName),
		Username:  NullDotStringToPointerString(u.FirstName),
		Email:     NullDotStringToPointerString(u.FirstName),
		Mobile:    NullDotStringToPointerString(u.FirstName),
		Phone:     NullDotStringToPointerString(u.FirstName),
		Address:   NullDotStringToPointerString(u.Address),
		Active:    NullDotBoolToPointerBool(u.Active),
	}
}
