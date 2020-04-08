package mock

import (
	"github.com/wednesday-solution/go-boiler"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(gorsk.User) (string, error)
}

// GenerateToken mock
func (j JWT) GenerateToken(u gorsk.User) (string, error) {
	return j.GenerateTokenFn(u)
}
