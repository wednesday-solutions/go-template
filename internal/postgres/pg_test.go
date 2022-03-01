package postgres_test

import (
	"testing"

	"go-template/internal/postgres"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	db, err := postgres.Connect()
	if err != nil {
		assert.NotNil(t, db)
	}
	db.Close()
}
