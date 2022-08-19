package cnvrttogql

import (
	graphql "go-template/gqlmodels"
	"go-template/models"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestUsersToGraphQlUsers(t *testing.T) {
	type args struct {
		u models.UserSlice
	}
	tests := []struct {
		name string
		args args
		want []*graphql.User
	}{
		{
			name: "Success",
			args: args{
				u: models.UserSlice{{
					ID: 1,
				}},
			},
			want: []*graphql.User{
				{
					ID: "1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UsersToGraphQlUsers(tt.args.u, 1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UsersToGraphQlUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleToGraphqlRole(t *testing.T) {
	type args struct {
		u *models.Role
	}
	tests := []struct {
		name string
		args args
		want *graphql.Role
	}{
		{
			name: "Success",
			args: args{
				u: &models.Role{
					ID: 1,
				},
			},
			want: &graphql.Role{
				ID: "1",
			},
		},
	}

	db, _, err := sqlmock.New()
	if err != nil {
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
	}
	boil.SetDB(db)
	defer db.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoleToGraphqlRole(tt.args.u, 1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoleToGraphqlRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
