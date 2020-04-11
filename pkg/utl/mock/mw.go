package mock

import (
	"github.com/wednesday-solutions/go-boiler"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(goboiler.User) (string, error)
}

// GenerateToken mock
func (j JWT) GenerateToken(u goboiler.User) (string, error) {
	return j.GenerateTokenFn(u)
}
