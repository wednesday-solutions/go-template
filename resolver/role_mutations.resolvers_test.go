package resolver_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/boil"
	fm "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/resolver"
)

func TestCreateRole(t *testing.T) {
	cases := []struct {
		name     string
		req      fm.RoleCreateInput
		wantResp *fm.RolePayload
		wantErr  bool
	}{
		{
			name:     "Fail on Create role",
			req:      fm.RoleCreateInput{},
			wantResp: &fm.RolePayload{},
			wantErr:  true,
		},
		{
			name:     "Success",
			req:      fm.RoleCreateInput{Name: "Role", AccessLevel: 100},
			wantResp: &fm.RolePayload{},
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
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

			if tt.name == "Fail on Create role" {
				// insert new user
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"roles\" (\"access_level\",\"name\",\"created_at\",\"updated_at\",\"deleted_at\") VALUES ($1,$2,$3,$4,$5)")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// insert new user
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"roles\" (\"access_level\",\"name\",\"created_at\",\"updated_at\",\"deleted_at\") VALUES ($1,$2,$3,$4,$5)")).
				WithArgs().
				WillReturnRows(rows)

			c := context.Background()
			response, _ := resolver1.Mutation().CreateRole(c, tt.req)
			assert.Equal(t, tt.wantResp, response)
		})
	}
}
