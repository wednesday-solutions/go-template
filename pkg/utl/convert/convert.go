package convert

import (
	"github.com/volatiletech/null"
	graphql "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/models"
	"strconv"
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
func UsersToGraphQlUsers(u models.UserSlice) []*graphql.User {
	var r []*graphql.User
	for _, e := range u {
		r = append(r, UserToGraphQlUser(e))
	}
	return r
}

// UserToGraphQlUser converts type models.User into pointer type graphql.User
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
