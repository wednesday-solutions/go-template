package mock

import (
	"github.com/wednesday-solutions/go-template/models"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(models.User) (string, error)
}

// GenerateToken mock
func (j JWT) GenerateToken(u models.User) (string, error) {
	return j.GenerateTokenFn(u)
}
