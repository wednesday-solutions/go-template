package resolver_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	fm "go-template/gqlmodels"
	"go-template/pkg/utl/convert"
	"go-template/resolver"
	"go-template/testutls"
	"regexp"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestCreateUser(t *testing.T) {
	cases := []struct {
		name     string
		req      fm.UserCreateInput
		wantResp *fm.User
		wantErr  bool
	}{
		{
			name:    "Fail on Create User",
			req:     fm.UserCreateInput{},
			wantErr: true,
		},
		{
			name: "Success",
			req: fm.UserCreateInput{
				FirstName: testutls.MockUser().FirstName.String,
				LastName:  testutls.MockUser().LastName.String,
				Username:  testutls.MockUser().Username.String,
				Email:     testutls.MockUser().Email.String,
				RoleID: fmt.Sprint(
					testutls.MockUser().RoleID.Int,
				),
			},
			wantResp: &fm.User{
				ID:                 fmt.Sprint(testutls.MockUser().ID),
				Email:              convert.NullDotStringToPointerString(testutls.MockUser().Email),
				FirstName:          convert.NullDotStringToPointerString(testutls.MockUser().FirstName),
				LastName:           convert.NullDotStringToPointerString(testutls.MockUser().LastName),
				Username:           convert.NullDotStringToPointerString(testutls.MockUser().Username),
				LastLogin:          convert.NullDotTimeToPointerInt(testutls.MockUser().LastLogin),
				LastPasswordChange: convert.NullDotTimeToPointerInt(testutls.MockUser().LastPasswordChange),
				DeletedAt:          convert.NullDotTimeToPointerInt(testutls.MockUser().DeletedAt),
				UpdatedAt:          convert.NullDotTimeToPointerInt(testutls.MockUser().UpdatedAt),
			},
			wantErr: false,
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				mock, db, _ := testutls.SetupEnvAndDB(
					t,
					testutls.Parameters{
						EnvFileLocation: "../.env.local",
					},
				)
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)

				if tt.name == "Fail on Create User" {
					// insert new user
					mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
						WithArgs().
						WillReturnError(fmt.Errorf(""))
				}

				ApplyFunc(
					bcrypt.GenerateFromPassword,
					func([]uint8, int) ([]uint8, error) {
						var a []uint8
						return a, nil
					},
				)
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`first_name`,`last_name`,`username`,`password`,`email`,"+
					"`mobile`,`address`,`active`,`last_login`,`last_password_change`,`token`,`role_id`,`created_at`,"+
					"`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(
						testutls.MockUser().FirstName,
						testutls.MockUser().LastName,
						testutls.MockUser().Username,
						"",
						testutls.MockUser().Email,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						testutls.MockUser().RoleID,
						testutls.AnyTime{},
						testutls.AnyTime{},
						nil,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				c := context.Background()
				response, err := resolver1.Mutation().CreateUser(c, tt.req)
				if tt.wantResp != nil {
					assert.EqualValues(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
			},
		)
	}
}

func TestUpdateUser(
	t *testing.T,
) {
	cases := []struct {
		name     string
		req      *fm.UserUpdateInput
		wantResp *fm.User
		wantErr  bool
	}{
		{
			name:    "Fail on finding User",
			req:     &fm.UserUpdateInput{},
			wantErr: true,
		},
		{
			name: "Success",
			req: &fm.UserUpdateInput{
				FirstName: &testutls.MockUser().FirstName.String,
				LastName:  &testutls.MockUser().LastName.String,
				Address:   &testutls.MockUser().Address.String,
			},
			wantResp: &fm.User{
				ID:        "0",
				FirstName: &testutls.MockUser().FirstName.String,
				LastName:  &testutls.MockUser().LastName.String,
				Address:   &testutls.MockUser().Address.String,
			},
			wantErr: false,
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				mock, db, _ := testutls.SetupEnvAndDB(
					t,
					testutls.Parameters{
						EnvFileLocation: "../.env.local",
					},
				)
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)

				if tt.wantErr {
					mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "users"`)).
						WithArgs().
						WillReturnError(fmt.Errorf(""))
				}

				rows := sqlmock.NewRows([]string{
					"first_name",
				}).
					AddRow(
						testutls.MockUser().
							FirstName,
					)
				mock.ExpectQuery(regexp.QuoteMeta("select * from `users`")).
					WithArgs(0).
					WillReturnRows(rows)
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `first_name`=?,`last_name`=?,`username`=?,`password`=?," +
					"`email`=?,`mobile`=?,`address`=?,`active`=?,`last_login`=?,`last_password_change`=?,`token`=?," +
					"`role_id`=?,`updated_at`=?,`deleted_at`=? WHERE `id`=?")).
					WillReturnResult(sqlmock.NewResult(1, 1))

				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				response, err := resolver1.Mutation().UpdateUser(ctx, tt.req)
				if tt.wantResp != nil &&
					response != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(
					t,
					tt.wantErr,
					err != nil,
				)
			},
		)
	}
}

func TestDeleteUser(
	t *testing.T,
) {
	cases := []struct {
		name     string
		wantResp *fm.UserDeletePayload
		wantErr  bool
	}{
		// {
		// 	name:    "Fail on finding user",
		// 	wantErr: true,
		// },
		{
			name: "Success",
			wantResp: &fm.UserDeletePayload{
				ID: "0",
			},
			wantErr: false,
		},
	}

	resolver1 := resolver.Resolver{}
	query := regexp.QuoteMeta("select * from `users` where `id`=?")
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				err := godotenv.Load("../.env.local")
				if err != nil {
					fmt.Print("error loading .env file")
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

				if tt.wantErr {
					mock.ExpectQuery(query).
						WithArgs().
						WillReturnError(fmt.Errorf(""))
				}
				// get user by id
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)
				mock.ExpectQuery(query).
					WithArgs().
					WillReturnRows(rows)
				// delete user
				result := driver.Result(
					driver.RowsAffected(
						1,
					),
				)
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE `id`=?")).
					WillReturnResult(result)

				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				response, err := resolver1.Mutation().
					DeleteUser(ctx)
				if tt.wantResp != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
			},
		)
	}
}
