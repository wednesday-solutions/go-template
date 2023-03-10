package cnvrttogql

import (
	graphql "go-template/gqlmodels"
	"go-template/models"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const SuccessCase = "Success"

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
			name: SuccessCase,
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
			name: SuccessCase,
			args: args{
				u: &models.Role{
					ID: 1,
				},
			},
			want: &graphql.Role{
				ID: "1",
			},
		},
		{
			name: SuccessCase,
			args: args{
				u: nil,
			},
			want: nil,
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

func TestUserToGraphQlUser(t *testing.T) {
	tests := []struct {
		name string
		req  *models.User
		want *graphql.User
	}{

		{
			name: SuccessCase,
			req:  nil,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserToGraphQlUser(tt.req, 0)
			assert.Equal(t, got, tt.want)

		})
	}
}
