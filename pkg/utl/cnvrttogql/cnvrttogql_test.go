package cnvrttogql

import (
	graphql "go-template/gqlmodels"
	"go-template/models"
	"reflect"
	"testing"
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
