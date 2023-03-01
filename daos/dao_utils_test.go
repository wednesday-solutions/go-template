package daos_test

import (
	"database/sql"
	"go-template/daos"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContextExecutor(t *testing.T) {

	cases := []struct {
		name string
		res  *sql.Tx
		err  error
	}{
		{
			name: "Passing role type value",
			res:  &sql.Tx{},
		},
	}

	// Loop through the test cases
	for _, tt := range cases {

		t.Run(tt.name, func(t *testing.T) {
			response := daos.GetContextExecutor(&sql.Tx{})

			// Check if the response is equal to the expected value
			assert.Equal(t, response, tt.res)

		})
	}
}
