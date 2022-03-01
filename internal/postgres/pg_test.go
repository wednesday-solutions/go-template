package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go-template/internal/postgres"
)

func TestConnect(t *testing.T) {
	db, err := postgres.Connect()
	if err != nil {
		assert.NotNil(t, db)
	}
	db.Close()
}
