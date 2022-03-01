package daos_test

import (
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/testutls"
	"regexp"
	"testing"
)

func TestCreateUserTx(t *testing.T) {

	cases := []struct {
		name string
		req  models.User
		err  error
	}{
		{
			name: "Passing user type value",
			req: models.User{
				ID:        0,
				FirstName: null.StringFrom("Wednesday"),
				LastName:  null.StringFrom("Solutions"),
				Username:  null.StringFrom("wedensday-solutions"),
				Password:  null.StringFrom("Abcd!1234"),
				Email:     null.StringFrom("hello@wednesday.is"),
				Mobile:    null.StringFrom("1000010000"),
				Phone:     null.StringFrom("022100010"),
				Address:   null.StringFrom("509, 510, Clover Hills Plaza, Khondwa"),
				Active:    null.Bool{bool(true), true},
				CreatedAt: null.Time{},
				UpdatedAt: null.Time{},
			},
			err: nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load("../.env.local")
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

		rows := sqlmock.NewRows([]string{
			"id",
			//"first_name",
			//"last_name",
			//"username",
			//"password",
			//"email",
			//"mobile",
			//"phone",
			//"address",
			//"active",
			//"created_at",
			//"updated_at",
			//"token",
			"last_login",
			"last_password_change",
			"token",
			"role_id",
			"deleted_at",
		}).AddRow(
			1,
			//null.StringFrom(tt.req.FirstName.String),
			//null.StringFrom(tt.req.LastName.String),
			//null.StringFrom(tt.req.Username.String),
			//null.StringFrom(tt.req.Password.String),
			//null.StringFrom(tt.req.Email.String),
			//null.StringFrom(tt.req.Mobile.String),
			//null.StringFrom(tt.req.Phone.String),
			//null.StringFrom(tt.req.Address.String),
			//null.Bool{tt.req.Active.Bool, true},
			null.Time{},
			null.Time{},
			null.String{},
			null.String{},
			null.Time{},
		)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("first_name","last_name","username","password","email","mobile","phone","address","active","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`)).
			WithArgs().WillReturnRows(rows)
		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.CreateUserTx(tt.req, nil)
			fmt.Println(err)
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
		err := godotenv.Load("../.env.local")
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
		err := godotenv.Load("../.env.local")
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
		err := godotenv.Load("../.env.local")
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
			req:  args{Token: testutls.MockToken},
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load("../.env.local")
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
		err := godotenv.Load("../.env.local")
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
		err := godotenv.Load("../.env.local")
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
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM \"users\" WHERE \"id\"=$1")).
			WillReturnResult(result)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.DeleteUser(tt.req)
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestFindAllUsersWithCount(t *testing.T) {

	oldDB := boil.GetDB()
	mock, db, _ := testutls.SetupEnvAndDB(t, testutls.Parameters{})

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
					Query: "SELECT * FROM \"users\";",
					DbResponse: sqlmock.NewRows([]string{"id", "email", "token"}).AddRow(
						testutls.MockID,
						testutls.MockEmail,
						testutls.MockToken),
				},
				{
					Query:      "SELECT COUNT(*) FROM \"users\";",
					DbResponse: sqlmock.NewRows([]string{"count"}).AddRow(testutls.MockCount),
				},
			},
		},
	}

	for _, tt := range cases {

		if tt.err != nil {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\";")).
				WithArgs().
				WillReturnError(fmt.Errorf("this is some error"))
		}

		for _, dbQuery := range tt.dbQueries {
			mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
				WithArgs().
				WillReturnRows(dbQuery.DbResponse)
		}

		t.Run(tt.name, func(t *testing.T) {
			res, c, err := daos.FindAllUsersWithCount([]qm.QueryMod{})
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
