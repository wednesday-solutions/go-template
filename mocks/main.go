package mocks

import (
	"github.com/volatiletech/null"
	"github.com/wednesday-solutions/go-template/models"
)

type key string

var (
	UserKey key = "user"
)

func MockUser() *models.User {
	return &models.User{
		FirstName: null.StringFrom("First"),
		LastName:  null.StringFrom("Last"),
		Username:  null.StringFrom("username"),
		Email:     null.StringFrom("mac@wednesday.is"),
		Mobile:    null.StringFrom("+911234567890"),
		Phone:     null.StringFrom("05943-1123"),
		Address:   null.StringFrom("22 Jump Street"),
	}
}
func MockUsers() []*models.User {
	return []*models.User{
		{
			FirstName: null.StringFrom("First"),
			LastName:  null.StringFrom("Last"),
			Username:  null.StringFrom("username"),
			Email:     null.StringFrom("mac@wednesday.is"),
			Mobile:    null.StringFrom("+911234567890"),
			Phone:     null.StringFrom("05943-1123"),
			Address:   null.StringFrom("22 Jump Street"),
		},
	}

}
