package query

import (
	"github.com/labstack/echo"

	"github.com/wednesday-solution/go-boiler"
)

// List prepares data for list queries
func List(u goboiler.AuthUser) (*goboiler.ListQuery, error) {
	switch true {
	case u.Role <= goboiler.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == goboiler.CompanyAdminRole:
		return &goboiler.ListQuery{Query: "company_id = ?", ID: u.CompanyID}, nil
	case u.Role == goboiler.LocationAdminRole:
		return &goboiler.ListQuery{Query: "location_id = ?", ID: u.LocationID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
