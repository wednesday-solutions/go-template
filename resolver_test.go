package goboiler_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	goboiler "github.com/wednesday-solutions/go-boiler"
	fm "github.com/wednesday-solutions/go-boiler/graphql_models"
	"github.com/wednesday-solutions/go-boiler/models"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/convert"
	"os"
	"regexp"
	"testing"
)

type key string

var (
	userKey key = "user"
)

func TestLogin(t *testing.T) {
	type args struct {
		UserName string
		Password string
	}
	cases := []struct {
		name     string
		req      args
		wantResp *fm.LoginResponse
		wantErr  bool
	}{
		{
			name:    "Fail on FindByUser",
			req:     args{UserName: "wednesday", Password: "pass123"},
			wantErr: true,
		},
		{
			name:     "Success",
			req:      args{UserName: "mac@wednesday.is", Password: "adminuser"},
			wantResp: &fm.LoginResponse{Token: "jwttokenstring", RefreshToken: "refreshtoken"},
		},
	}

	resolver := goboiler.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
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

			if tt.name == "Fail on FindByUser" {
				// get user by username
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (username=$1) LIMIT 1;")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// get user by username
			rows := sqlmock.NewRows([]string{"id", "password", "active"}).AddRow(1, "$2a$10$dS5vK8hHmG5gzwV8f7TK5.WHviMBqmYQLYp30a3XvqhCW9Wvl2tOS", true)
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (username=$1) LIMIT 1;")).
				WithArgs().
				WillReturnRows(rows)

			// update users with token
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta("UPDATE \"users\" ")).
				WillReturnResult(result)

			c := context.Background()
			response, err := resolver.Mutation().Login(c, tt.req.UserName, tt.req.Password)
			if tt.wantResp != nil && response != nil {
				tt.wantResp.RefreshToken = response.RefreshToken
				tt.wantResp.Token = response.Token
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestMe(t *testing.T) {
	cases := []struct {
		name     string
		wantResp *fm.User
		wantErr  bool
	}{
		{
			name:     "Success",
			wantErr:  false,
			wantResp: &fm.User{FirstName: convert.StringToPointerString("First"), LastName: convert.StringToPointerString("Last"), Username: convert.StringToPointerString("username"), Email: convert.StringToPointerString("mac@wednesday.is"), Mobile: convert.StringToPointerString("+911234567890"), Phone: convert.StringToPointerString("05943-1123"), Address: convert.StringToPointerString("22 Jump Street")},
		},
	}

	resolver := goboiler.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
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

			// get user by id
			rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "mobile", "username", "phone", "address"}).AddRow(1, "mac@wednesday.is", "First", "Last", "+911234567890", "username", "05943-1123", "22 Jump Street")
			mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
				WithArgs().
				WillReturnRows(rows)

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
			response, err := resolver.Query().Me(ctx)
			assert.Equal(t, tt.wantResp, response)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUsers(t *testing.T) {
	cases := []struct {
		name       string
		pagination *fm.UserPagination
		wantResp   []*models.User
		wantErr    bool
	}{
		{
			name:    "Success",
			wantErr: false,
			wantResp: []*models.User{
				{
					FirstName: null.StringFrom("First"),
					LastName:  null.StringFrom("Last"),
					Username:  null.StringFrom("username"),
					Email:     null.StringFrom("mac@wednesday.is"),
					Mobile:    null.StringFrom("+911234567890"),
					Phone:     null.StringFrom("05943-1123"),
					Address:   null.StringFrom("22 Jump Street"),
				},
			},
		},
	}

	resolver := goboiler.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
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

			// get user by id
			rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "mobile", "username", "phone", "address"}).AddRow(1, "mac@wednesday.is", "First", "Last", "+911234567890", "username", "05943-1123", "22 Jump Street")
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\";")).
				WithArgs().
				WillReturnRows(rows)

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
			response, err := resolver.Query().Users(ctx, tt.pagination)
			if tt.wantResp != nil && response != nil {
				assert.Equal(t, len(tt.wantResp), len(response.Users))
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestChangePassword(t *testing.T) {

	type changeReq struct {
		OldPassword string
		NewPassword string
	}
	cases := []struct {
		name     string
		req      changeReq
		wantResp *fm.ChangePasswordResponse
		wantErr  bool
	}{
		{
			name:    "Fail on FindByUser",
			req:     changeReq{OldPassword: "adminuser!A9@@@@", NewPassword: "adminuser!A9@"},
			wantErr: true,
		},
		{
			name:    "Incorrect Old Password",
			req:     changeReq{OldPassword: "admin", NewPassword: "adminuser!A9@"},
			wantErr: true,
		},
		{
			name:     "Success",
			req:      changeReq{OldPassword: "adminuser", NewPassword: "adminuser!A9@"},
			wantResp: &fm.ChangePasswordResponse{Ok: true},
			wantErr:  false,
		},
	}

	resolver := goboiler.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
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

			if tt.name == "Fail on FindByUser" {
				// get user by id
				mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// get user by id
			rows := sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(1, "mac@wednesday.is", "$2a$10$dS5vK8hHmG5gzwV8f7TK5.WHviMBqmYQLYp30a3XvqhCW9Wvl2tOS")
			mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
				WithArgs().
				WillReturnRows(rows)
			if tt.name == "Success" {
				// update password
				result := driver.Result(driver.RowsAffected(1))
				mock.ExpectExec(regexp.QuoteMeta("UPDATE \"users\" ")).
					WillReturnResult(result)
			}

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
			response, err := resolver.Mutation().ChangePassword(ctx, tt.req.OldPassword, tt.req.NewPassword)
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		wantResp *fm.RefreshTokenResponse
		wantErr  bool
	}{
		{
			name:    "Fail on FindByToken",
			req:     "refreshToken",
			wantErr: true,
		},
		{
			name:     "Success",
			req:      "refresh_token",
			wantResp: &fm.RefreshTokenResponse{Token: "token_string"},
			wantErr:  false,
		},
	}

	resolver := goboiler.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
			if err != nil {
				fmt.Print("Error loading .env file")
			}
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			oldDB := boil.GetDB()
			defer func() {
				db.Close()
				boil.SetDB(oldDB)
			}()
			boil.SetDB(db)

			if tt.name == "Fail on FindByToken" {
				// get user by token
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (token=$1) LIMIT 1;")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// get user by token
			rows := sqlmock.NewRows([]string{"id", "email", "token"}).AddRow(1, "mac@wednesday.is", "token_string")
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (token=$1) LIMIT 1;")).
				WithArgs().
				WillReturnRows(rows)

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
			response, err := resolver.Mutation().RefreshToken(ctx, tt.req)
			if tt.wantResp != nil && response != nil {
				tt.wantResp.Token = response.Token
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestCreateUser(t *testing.T) {
	cases := []struct {
		name     string
		req      fm.UserCreateInput
		wantResp *fm.UserPayload
		wantErr  bool
	}{
		{
			name:    "Fail on Create User",
			req:     fm.UserCreateInput{},
			wantErr: true,
		},
		{
			name:     "Success",
			req:      fm.UserCreateInput{FirstName: convert.StringToPointerString("First"), LastName: convert.StringToPointerString("Last"), Username: convert.StringToPointerString("username"), Email: convert.StringToPointerString("mac@wednesday.is")},
			wantResp: &fm.UserPayload{User: &fm.User{FirstName: convert.StringToPointerString("First"), LastName: convert.StringToPointerString("Last"), Username: convert.StringToPointerString("username"), Email: convert.StringToPointerString("mac@wednesday.is")}},
			wantErr:  false,
		},
	}

	resolver := goboiler.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
			if err != nil {
				fmt.Print("Error loading .env file")
			}
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			oldDB := boil.GetDB()
			defer func() {
				db.Close()
				boil.SetDB(oldDB)
			}()
			boil.SetDB(db)

			if tt.name == "Fail on Create User" {
				// insert new user
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"first_name\",\"last_name\",\"username\",\"password\",\"email\",\"mobile\",\"phone\",\"address\",\"active\",\"last_login\",\"last_password_change\",\"token\",\"role_id\",\"company_id\",\"location_id\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// insert new user
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"first_name\",\"last_name\",\"username\",\"password\",\"email\",\"mobile\",\"phone\",\"address\",\"active\",\"last_login\",\"last_password_change\",\"token\",\"role_id\",\"company_id\",\"location_id\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)")).
				WithArgs().
				WillReturnRows(rows)

			c := context.Background()
			response, err := resolver.Mutation().CreateUser(c, tt.req)
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	cases := []struct {
		name     string
		req      *fm.UserUpdateInput
		wantResp *fm.UserUpdatePayload
		wantErr  bool
	}{
		{
			name:    "Fail on finding User",
			req:     &fm.UserUpdateInput{},
			wantErr: true,
		},
		{
			name:     "Success",
			req:      &fm.UserUpdateInput{FirstName: convert.StringToPointerString("First"), LastName: convert.StringToPointerString("Last"), Address: convert.StringToPointerString("address")},
			wantResp: &fm.UserUpdatePayload{Ok: true},
			wantErr:  false,
		},
	}

	resolver := goboiler.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
			if err != nil {
				fmt.Print("Error loading .env file")
			}
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			oldDB := boil.GetDB()
			defer func() {
				db.Close()
				boil.SetDB(oldDB)
			}()
			boil.SetDB(db)

			if tt.name == "Fail on finding User" {
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE \"users\" SET \"updated_at\"=$1,\"deleted_at\"=$2,\"first_name\"=$3,\"last_name\"=$4,\"username\"=$5,\"password\"=$6,\"email\"=$7,\"mobile\"=$8,\"phone\"=$9,\"address\"=$10,\"active\"=$11,\"last_login\"=$12,\"last_password_change\"=$13,\"token\"=$14,\"role_id\"=$15,\"company_id\"=$16,\"location_id\"=$17 WHERE \"id\"=$18")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// update users with new information
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta("UPDATE \"users\" SET \"updated_at\"=$1,\"deleted_at\"=$2,\"first_name\"=$3,\"last_name\"=$4,\"username\"=$5,\"password\"=$6,\"email\"=$7,\"mobile\"=$8,\"phone\"=$9,\"address\"=$10,\"active\"=$11,\"last_login\"=$12,\"last_password_change\"=$13,\"token\"=$14,\"role_id\"=$15,\"company_id\"=$16,\"location_id\"=$17 WHERE \"id\"=$18")).
				WillReturnResult(result)

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
			response, err := resolver.Mutation().UpdateUser(ctx, tt.req)
			if tt.wantResp != nil && response != nil {
				assert.Equal(t, tt.wantResp.Ok, response.Ok)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	cases := []struct {
		name     string
		wantResp *fm.UserDeletePayload
		wantErr  bool
	}{
		{
			name:    "Fail on finding user",
			wantErr: true,
		},
		{
			name:     "Success",
			wantResp: &fm.UserDeletePayload{ID: "0"},
			wantErr:  false,
		},
	}

	resolver := goboiler.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
			if err != nil {
				fmt.Print("Error loading .env file")
			}
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			oldDB := boil.GetDB()
			defer func() {
				db.Close()
				boil.SetDB(oldDB)
			}()
			boil.SetDB(db)

			if tt.name == "Fail on finding user" {
				mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// get user by id
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
				WithArgs().
				WillReturnRows(rows)
			// delete user
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM \"users\" WHERE \"id\"=$1")).
				WillReturnResult(result)

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
			response, err := resolver.Mutation().DeleteUser(ctx)
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
