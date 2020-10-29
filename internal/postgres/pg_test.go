package postgres_test

import (
	"database/sql"
	"testing"

	"github.com/fortytw2/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/wednesday-solutions/go-template/internal/postgres"
)

func TestNew(t *testing.T) {
	container, err := dockertest.RunContainer("postgres:alpine", "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			return err
		}

		return db.Ping()
	}, "-e", "POSTGRES_PASSWORD=postgres", "-e", "POSTGRES_USER=postgres")
	defer container.Shutdown()
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	assert.NotNil(t, container)
}

func TestConnect(t *testing.T) {
	db, err := postgres.Connect()
	if err != nil {
		assert.NotNil(t, db)
	}
	db.Close()
}
