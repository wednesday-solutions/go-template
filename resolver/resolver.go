package resolver

import (
	"sync"

	fm "github.com/wednesday-solutions/go-template/graphql_models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver ...
type Resolver struct {
	sync.Mutex
	Observers map[string]chan *fm.User
}
