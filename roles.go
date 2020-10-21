package gotemplate

// AccessRole represents access role type
type AccessRole int

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = 200

	// AdminRole has admin specific permissions
	AdminRole AccessRole = 150

	// UserRole is a standard user
	UserRole AccessRole = 100
)
