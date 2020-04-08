package query_test

import (
	"testing"

	"github.com/labstack/echo"

	"github.com/wednesday-solution/go-boiler"

	"github.com/stretchr/testify/assert"

	"github.com/wednesday-solution/go-boiler/pkg/utl/query"
)

func TestList(t *testing.T) {
	type args struct {
		user goboiler.AuthUser
	}
	cases := []struct {
		name     string
		args     args
		wantData *goboiler.ListQuery
		wantErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: goboiler.AuthUser{
				Role: goboiler.SuperAdminRole,
			}},
		},
		{
			name: "Company admin user",
			args: args{user: goboiler.AuthUser{
				Role:      goboiler.CompanyAdminRole,
				CompanyID: 1,
			}},
			wantData: &goboiler.ListQuery{
				Query: "company_id = ?",
				ID:    1},
		},
		{
			name: "Location admin user",
			args: args{user: goboiler.AuthUser{
				Role:       goboiler.LocationAdminRole,
				CompanyID:  1,
				LocationID: 2,
			}},
			wantData: &goboiler.ListQuery{
				Query: "location_id = ?",
				ID:    2},
		},
		{
			name: "Normal user",
			args: args{user: goboiler.AuthUser{
				Role: goboiler.UserRole,
			}},
			wantErr: echo.ErrForbidden,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.List(tt.args.user)
			assert.Equal(t, tt.wantData, q)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
