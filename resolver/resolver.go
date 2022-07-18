package resolver

import (
	"sync"

	fm "go-template/gqlmodels"
)

// This file will
// not be
// regenerated
// automatically.
//
// It serves as
// dependency
// injection for
// your app, add any
// dependencies you
// require here.

// Resolver ...
type Resolver struct {
	sync.Mutex
	Observers map[string]chan *fm.User
}
