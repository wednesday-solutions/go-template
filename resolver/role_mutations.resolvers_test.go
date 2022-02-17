package resolver_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
	fm "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/resolver"
)

func TestCreateRole(t *testing.T) {
	cases := []struct {
		name     string
		req      fm.RoleCreateInput
		wantResp *fm.Role
		wantErr  bool
	}{
		{
			name:     "Fail on Create role",
			req:      fm.RoleCreateInput{},
			wantResp: &fm.Role{},
			wantErr:  true,
		},
		{
			name:     "Success",
			req:      fm.RoleCreateInput{Name: "Role", AccessLevel: 200},
			wantResp: &fm.Role{},
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
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
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

			mock.ExpectQuery(regexp.QuoteMeta("select * from \"roles\" where \"id\"=$1")).
				WithArgs([]driver.Value{0}...).
				WillReturnRows(rows)

			mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
				WithArgs([]driver.Value{0}...).
				WillReturnRows(rows)

			if tt.name == "Fail on Create role" {
				// insert new user
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"roles\" (\"access_level\"," +
					"\"name\",\"created_at\",\"updated_at\",\"deleted_at\") VALUES ($1,$2,$3,$4,$5)")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// insert new user
			rows = sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"roles\" (\"access_level\"," +
				"\"name\",\"created_at\",\"updated_at\",\"deleted_at\") VALUES ($1,$2,$3,$4,$5)")).
				WithArgs().
				WillReturnRows(rows)

			c := context.Background()
			response, err := resolver1.Mutation().CreateRole(c, tt.req)
			if tt.wantErr {
				assert.NotNil(t, err)
			}
			assert.Equal(t, tt.wantResp, response)
		})
	}
}
