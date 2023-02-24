package constants

// AccessRole represents access role type
type AccessRole int

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = 100

	// AdminRole has admin specific permissions
	AdminRole AccessRole = 110

	// UserRole is a standard user
	UserRole AccessRole = 200

	// CompanyAdmin has admin specific permissions
	COMPANY_ADMIN AccessRole = 120

	// LocationAdmin has admin specific permissions
	LOCATION_ADMIN AccessRole = 130
)

const (
	MaxDepth = 4
)
