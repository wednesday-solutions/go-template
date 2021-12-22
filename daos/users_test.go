package daos_test

import (
	"database/sql/driver"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/models"
)

func TestCreateUserTx(t *testing.T) {

	cases := []struct {
		name string
		req  models.User
		err  error
	}{
		{
			name: "Passing user type value",
			req:  models.User{},
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
		if err != nil {
			fmt.Print("Error loading .env file")
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// Inject mock instance into boil.
		oldDB := boil.GetDB()
		defer func() {
			db.Close()
			boil.SetDB(oldDB)
		}()
		boil.SetDB(db)

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"users\" (\"first_name\",\"last_name\"," +
			"\"username\",\"password\",\"email\",\"mobile\",\"phone\",\"address\",\"active\",\"last_login\"," +
			"\"last_password_change\",\"token\",\"role_id\",\"created_at\",\"updated_at\",\"deleted_at\") " +
			"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)")).
			WithArgs().
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.CreateUserTx(tt.req, nil)
			if err != nil {
				assert.Equal(t, true, tt.err != nil)
			} else {
				assert.Equal(t, err, tt.err)
			}
		})
	}
}

func TestFindUserByID(t *testing.T) {

	cases := []struct {
		name string
		req  int
		err  error
	}{
		{
			name: "Passing a user_id",
			req:  1,
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
		if err != nil {
			fmt.Print("Error loading .env file")
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// Inject mock instance into boil.
		oldDB := boil.GetDB()
		defer func() {
			db.Close()
			boil.SetDB(oldDB)
		}()
		boil.SetDB(db)

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
			WithArgs().
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.FindUserByID(tt.req)
			assert.Equal(t, err, tt.err)

		})
	}
}

func TestFindUserByEmail(t *testing.T) {

	type args struct {
		email string
	}
	cases := []struct {
		name string
		req  args
		err  error
	}{
		{
			name: "Fail on finding user",
			req:  args{email: "abc"},
			err:  fmt.Errorf("sql: no rows in sql"),
		},
		{
			name: "Passing an email",
			req:  args{email: "mac"},
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
		if err != nil {
			fmt.Print("Error loading .env file")
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// Inject mock instance into boil.
		oldDB := boil.GetDB()
		defer func() {
			db.Close()
			boil.SetDB(oldDB)
		}()
		boil.SetDB(db)

		if tt.name == "Fail on finding user" {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (email=$1) LIMIT 1;")).
				WithArgs().
				WillReturnError(fmt.Errorf(""))
		}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (email=$1) LIMIT 1;")).
			WithArgs().
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.FindUserByEmail(tt.req.email)
			if err != nil {
				assert.Equal(t, true, tt.err != nil)
			} else {
				assert.Equal(t, err, tt.err)
			}
		})
	}
}

func TestFindUserByUserName(t *testing.T) {

	type args struct {
		Username string
	}
	cases := []struct {
		name string
		req  args
		err  error
	}{
		{
			name: "Fail on finding user username",
			req:  args{Username: "user"},
			err:  fmt.Errorf("sql: no rows in sql"),
		},
		{
			name: "Passing a valid username",
			req:  args{Username: "user_name"},
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
		if err != nil {
			fmt.Print("Error loading .env file")
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// Inject mock instance into boil.
		oldDB := boil.GetDB()
		defer func() {
			db.Close()
			boil.SetDB(oldDB)
		}()
		boil.SetDB(db)

		if tt.name == "Fail on finding user username" {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (username=$1) LIMIT 1;")).
				WithArgs().
				WillReturnError(fmt.Errorf(""))
		}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (username=$1) LIMIT 1;")).
			WithArgs().
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.FindUserByUserName(tt.req.Username)
			if err != nil {
				assert.Equal(t, true, tt.err != nil)
			} else {
				assert.Equal(t, err, tt.err)
			}
		})
	}
}

func TestFindUserByToken(t *testing.T) {

	type args struct {
		Token string
	}
	cases := []struct {
		name string
		req  args
		err  error
	}{
		{
			name: "Fail on finding user token",
			req:  args{Token: "tokenString"},
			err:  fmt.Errorf("sql: no rows in sql"),
		},
		{
			name: "Passing an email",
			req:  args{Token: "token_string"},
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
		if err != nil {
			fmt.Print("Error loading .env file")
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// Inject mock instance into boil.
		oldDB := boil.GetDB()
		defer func() {
			db.Close()
			boil.SetDB(oldDB)
		}()
		boil.SetDB(db)

		if tt.name == "Fail on finding user token" {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (token=$1) LIMIT 1;")).
				WithArgs().
				WillReturnError(fmt.Errorf(""))
		}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (token=$1) LIMIT 1;")).
			WithArgs().
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.FindUserByToken(tt.req.Token)
			if err != nil {
				assert.Equal(t, true, tt.err != nil)
			} else {
				assert.Equal(t, err, tt.err)
			}
		})
	}
}

func TestUpdateUserTx(t *testing.T) {

	cases := []struct {
		name string
		req  models.User
		err  error
	}{
		{
			name: "Passing user type value",
			req:  models.User{},
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
		if err != nil {
			fmt.Print("Error loading .env file")
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// Inject mock instance into boil.
		oldDB := boil.GetDB()
		defer func() {
			db.Close()
			boil.SetDB(oldDB)
		}()
		boil.SetDB(db)

		result := driver.Result(driver.RowsAffected(1))
		// get access_token
		mock.ExpectExec(regexp.QuoteMeta("UPDATE \"users\" ")).
			WillReturnResult(result)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.UpdateUserTx(tt.req, nil)
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestDeleteUser(t *testing.T) {

	cases := []struct {
		name string
		req  models.User
		err  error
	}{
		{
			name: "Passing user type value",
			req:  models.User{},
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
		if err != nil {
			fmt.Print("Error loading .env file")
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// Inject mock instance into boil.
		oldDB := boil.GetDB()
		defer func() {
			db.Close()
			boil.SetDB(oldDB)
		}()
		boil.SetDB(db)

		// delete user
		result := driver.Result(driver.RowsAffected(1))
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM \"users\" WHERE \"id\"=$1")).
			WillReturnResult(result)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.DeleteUser(tt.req)
			assert.Equal(t, err, tt.err)
		})
	}
}
