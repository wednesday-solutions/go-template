package resolver_test

import (
	"context"
	"errors"
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
			name:     "RedisCache Error",
			req:      fm.RoleCreateInput{},
			wantResp: &fm.RolePayload{},
			wantErr:  true,
		},
		{
			name:     "RedisCache GetRole Error",
			req:      fm.RoleCreateInput{},
			wantResp: &fm.RolePayload{},
			wantErr:  true,
		},
		{
			name: "Unauthorized",
			req: fm.RoleCreateInput{
				Name:        "Role",
				AccessLevel: 200,
			},
			wantResp: &fm.RolePayload{},
		},

		{
			name: "Success",
			req: fm.RoleCreateInput{
				Name:        "Role",
				AccessLevel: 200,
			},
			wantResp: &fm.RolePayload{Role: &fm.Role{

				AccessLevel: 200,
				Name:        "Role",
			}},
		},
		{
			name: "CreateRole Error",
			req: fm.RoleCreateInput{
				Name:        "UserRole",
				AccessLevel: 200,
			},
			wantErr: true,
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {

		patchUserID := gomonkey.ApplyFunc(auth.UserIDFromContext,
			func(ctx context.Context) int {
				return 1
			})
		patchGetUser := gomonkey.ApplyFunc(rediscache.GetUser,
			func(userID int, ctx context.Context) (*models.User, error) {
				return &models.User{
					RoleID: null.IntFrom(1),
				}, nil
			})
		patchGetRole := gomonkey.ApplyFunc(rediscache.GetRole,
			func(roleID int, ctx context.Context) (*models.Role, error) {
				return &models.Role{
					AccessLevel: 100,
					Name:        "SuperAdminRole",
				}, nil
			})
		patchCreateRole := gomonkey.ApplyFunc(daos.CreateRole,
			func(role models.Role, ctx context.Context) (models.Role, error) {
				return models.Role{
					AccessLevel: 200,
					Name:        "Role",
				}, nil
			})
		//defer patchDaos1.Reset()
		defer patchUserID.Reset()
		defer patchGetUser.Reset()
		defer patchGetRole.Reset()
		defer patchCreateRole.Reset()
		t.Run(tt.name,
			func(t *testing.T) {

				if tt.name == "RedisCache Error" {
					patchGetUser := gomonkey.ApplyFunc(rediscache.GetUser,
						func(userID int, ctx context.Context) (*models.User, error) {
							return nil, errors.New("redis cache")
						})
					defer patchGetUser.Reset()
				}

				if tt.name == "RedisCache GetRole Error" {
					patchGetRole := gomonkey.ApplyFunc(rediscache.GetRole,
						func(roleID int, ctx context.Context) (*models.Role, error) {
							return nil, errors.New("data")
						})
					defer patchGetRole.Reset()
				}
				if tt.name == "Unauthorized" {
					patchGetRole := gomonkey.ApplyFunc(rediscache.GetRole,
						func(roleID int, ctx context.Context) (*models.Role, error) {
							return &models.Role{
								AccessLevel: 200,
								Name:        "Role",
							}, nil
						})
					defer patchGetRole.Reset()
				}

				if tt.name == "CreateRole Error" {
					patchCreateRole := gomonkey.ApplyFunc(daos.CreateRole,
						func(role models.Role, ctx context.Context) (models.Role, error) {
							return models.Role{}, errors.New("error")
						})

					defer patchCreateRole.Reset()
				}
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

				if tt.wantErr {
					assert.NotNil(t, err)
				}
				log.Println(tt.wantResp)
				assert.Equal(t, tt.wantResp, response)

				//	},
				//)
			})
	}
}
