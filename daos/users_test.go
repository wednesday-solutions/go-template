package daos_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log"
	"regexp"
	"testing"

	"go-template/daos"
	"go-template/internal/config"
	"go-template/models"
	"go-template/pkg/utl/convert"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const ErrorFindingUser = "Fail on finding user"

func TestCreateUserTx(t *testing.T) {

	cases := []struct {
		name string
		req  models.User
		err  error
	}{
		{
			name: "Passing user type value",
			req: models.User{
				ID:       testutls.MockUser().ID,
				Email:    testutls.MockUser().Email,
				Password: testutls.MockUser().Password,
			},
			err: nil,
		},
	}

	for _, tt := range cases {
		// Inject mock instance into boil.

		err := config.LoadEnvWithFilePrefix(convert.StringToPointerString("./../"))
		if err != nil {
			log.Fatal(err)
		}
		mock, db, _ := testutls.SetupMockDB(t)
		oldDB := boil.GetDB()
		defer func() {
			db.Close()
			boil.SetDB(oldDB)
		}()
		boil.SetDB(db)

		rows := sqlmock.NewRows([]string{
			"first_name",
			"last_name",
			"username",
			"mobile",
			"address",
			"active",
			"last_login",
			"last_password_change",
			"token",
			"role_id",
			"deleted_at",
		}).AddRow(
			testutls.MockUser().FirstName,
			testutls.MockUser().LastName,
			testutls.MockUser().Username,
			testutls.MockUser().Mobile,
			testutls.MockUser().Address,
			testutls.MockUser().Active,
			testutls.MockUser().LastLogin,
			testutls.MockUser().LastPasswordChange,
			testutls.MockUser().Token,
			testutls.MockUser().RoleID,
			testutls.MockUser().DeletedAt,
		)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WithArgs().
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.CreateUser(tt.req, context.Background())
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
		err := config.LoadEnv()
		if err != nil {
			fmt.Print("error loading .env file")
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

		mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
			WithArgs().
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.FindUserByID(tt.req, context.Background())
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
			name: ErrorFindingUser,
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
		err := config.LoadEnv()
		if err != nil {
			fmt.Print("error loading .env file")
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

		if tt.name == ErrorFindingUser {
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (email=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnError(fmt.Errorf(""))
		}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (email=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.FindUserByEmail(tt.req.email, context.Background())
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
		err := config.LoadEnv()
		if err != nil {
			fmt.Print("error loading .env file")
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
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (username=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnError(fmt.Errorf(""))
		}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (username=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.FindUserByUserName(tt.req.Username, context.Background())
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
		err := config.LoadEnv()
		if err != nil {
			fmt.Print("error loading .env file")
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
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).
			WillReturnResult(result)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.UpdateUser(tt.req, context.Background())
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
		err := config.LoadEnv()
		if err != nil {
			fmt.Print("error loading .env file")
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
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "id"=$1`)).
			WillReturnResult(result)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.DeleteUser(tt.req, context.Background())
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestFindAllUsersWithCount(t *testing.T) {

	oldDB := boil.GetDB()
	err := config.LoadEnvWithFilePrefix(convert.StringToPointerString("./../"))
	if err != nil {
		log.Fatal(err)
	}
	mock, db, _ := testutls.SetupMockDB(t)

	cases := []struct {
		name      string
		err       error
		dbQueries []testutls.QueryData
	}{
		{
			name: "Failed to find all users with count",
			err:  fmt.Errorf("sql: no rows in sql"),
		},
		{
			name: "Successfully find all users with count",
			err:  nil,
			dbQueries: []testutls.QueryData{
				{
					Query: `SELECT "users".* FROM "users";`,
					DbResponse: sqlmock.NewRows([]string{"id", "email", "token"}).AddRow(
						testutls.MockID,
						testutls.MockEmail,
						testutls.MockToken),
				},
				{
					Query:      `SELECT COUNT(*) FROM "users";`,
					DbResponse: sqlmock.NewRows([]string{"count"}).AddRow(testutls.MockCount),
				},
			},
		},
	}

	for _, tt := range cases {

		if tt.err != nil {
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users";`)).
				WithArgs().
				WillReturnError(fmt.Errorf("this is some error"))
		}

		for _, dbQuery := range tt.dbQueries {
			mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
				WithArgs().
				WillReturnRows(dbQuery.DbResponse)
		}

		t.Run(tt.name, func(t *testing.T) {
			res, c, err := daos.FindAllUsersWithCount([]qm.QueryMod{}, context.Background())
			if err != nil {
				assert.Equal(t, true, tt.err != nil)
			} else {
				assert.Equal(t, err, tt.err)
				assert.Equal(t, testutls.MockCount, c)
				assert.Equal(t, res[0].Email, null.StringFrom(testutls.MockEmail))
				assert.Equal(t, res[0].Token, null.StringFrom(testutls.MockToken))
				assert.Equal(t, res[0].ID, int(testutls.MockID))

			}
		})
	}
	boil.SetDB(oldDB)
	db.Close()
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
			req:  args{Token: testutls.MockToken},
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := config.LoadEnv()
		if err != nil {
			fmt.Print("error loading .env file")
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
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnError(fmt.Errorf(""))
		}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.FindUserByToken(tt.req.Token, context.Background())
			if err != nil {
				assert.Equal(t, true, tt.err != nil)
			} else {
				assert.Equal(t, err, tt.err)
			}
		})
	}
}
