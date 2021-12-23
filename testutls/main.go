package testutls

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
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

func SetupEnvAndDB(t *testing.T) (mock sqlmock.Sqlmock, db *sql.DB, err error) {
	err = godotenv.Load("../.env.local")
	if err != nil {
		fmt.Print("error loading .env file")
	}
	db, mock, err = sqlmock.New()
	if err != nil {
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
	}
	boil.SetDB(db)
	return mock, db, nil
}
