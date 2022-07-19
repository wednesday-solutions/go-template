package mysql_test

import (
	"testing"

	"go-template/internal/mysql"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	db, err := mysql.Connect()
	if err != nil {
		assert.NotNil(t, db)
	}
	db.Close()
}
