package resolver_test

import (
	"context"
	"go-template/daos"
	"go-template/internal/middleware/auth"
	"go-template/models"
	"go-template/pkg/utl/rediscache"
	"go-template/resolver"
	"log"
	"testing"

	"github.com/volatiletech/null/v8"

	"github.com/agiledragon/gomonkey/v2"

	fm "go-template/gqlmodels"

	"github.com/stretchr/testify/assert"
)

func TestCreateRole(
	t *testing.T,
) {
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
			name: "Success",
			req: fm.RoleCreateInput{
				Name:        "Role",
				AccessLevel: 200,
			},
			wantResp: &fm.RolePayload{},
		},
		{
			name: "Fail on Create role",
			req: fm.RoleCreateInput{
				Name:        "UserRole",
				AccessLevel: 100,
			},
			wantResp: &fm.RolePayload{},
			wantErr:  true,
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {

		patchDaos2 := gomonkey.ApplyFunc(auth.UserIDFromContext,
			func(ctx context.Context) int {
				return 1
			})
		patchDaos3 := gomonkey.ApplyFunc(rediscache.GetUser,
			func(userID int, ctx context.Context) (*models.User, error) {
				return &models.User{
					RoleID: null.IntFrom(1),
				}, nil
			})
		patchDaos4 := gomonkey.ApplyFunc(rediscache.GetRole,
			func(roleID int, ctx context.Context) (*models.Role, error) {
				return &models.Role{
					AccessLevel: 100,
					Name:        "SuperAdminRole",
				}, nil
			})
		patchDaos5 := gomonkey.ApplyFunc(daos.CreateRole,
			func(role models.Role, ctx context.Context) (models.Role, error) {
				return models.Role{
					AccessLevel: 200,
					Name:        "Role",
				}, nil
			})
		//defer patchDaos1.Reset()
		defer patchDaos2.Reset()
		defer patchDaos3.Reset()
		defer patchDaos4.Reset()
		defer patchDaos5.Reset()
		t.Run(tt.name,
			func(t *testing.T) {
				//		err := godotenv.Load("../.env.local")
				//		if err != nil {
				//			fmt.Print("error loading .env file")
				//		}
				//		db, mock, err := sqlmock.New()
				//		if err != nil {
				//			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				//		}
				//		oldDB := boil.GetDB()
				//		defer func() {
				//			db.Close()
				//			boil.SetDB(oldDB)
				//		}()
				//		boil.SetDB(db)
				//		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				//
				//		mock.ExpectQuery(`select * from "roles" where "id"=$1`).
				//			WithArgs([]driver.Value{0}...).
				//			WillReturnRows(rows)
				//
				//		mock.ExpectQuery(`select * from "users" where "id"=$1`).
				//			WithArgs([]driver.Value{0}...).
				//			WillReturnRows(rows)
				//
				//		if tt.name == "Fail on Create role" {
				//			// insert new user
				//			mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "roles"`)).
				//				WithArgs().
				//				WillReturnError(fmt.Errorf(""))
				//		}
				//		// insert new user
				//		rows = sqlmock.NewRows([]string{"id"}).AddRow(1)
				//		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "roles"`)).
				//			WithArgs().
				//			WillReturnRows(rows)
				//
				c := context.Background()
				response, err := resolver1.Mutation().CreateRole(c, tt.req)
				log.Println(tt.wantResp)
				log.Println(response)
				if tt.wantErr {
					assert.NotNil(t, err)
				}
				assert.Equal(t, tt.wantResp, response)

				//	},
				//)
			})
	}
}
