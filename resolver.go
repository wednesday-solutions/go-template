// Package goboiler ...
// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.
package goboiler

import (
	"database/sql"
)

// Resolver ...
type Resolver struct {
	db *sql.DB
}

const inputKey = "input"

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
