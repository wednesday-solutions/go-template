package daos_test

import (
	"database/sql"
	"go-template/daos"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

	// Loop through the test cases and apply gomonkey patches as required to mock the function call
	for _, tt := range cases {

		patchDaos := gomonkey.ApplyFunc(daos.GetContextExecutor,
			func(tx *sql.Tx) (contextExecutor boil.ContextExecutor) {
				return boil.ContextExecutor(tx)
			})

		// Defer resetting of the monkey patches.
		defer patchDaos.Reset()
		t.Run(tt.name, func(t *testing.T) {
			response := daos.GetContextExecutor(&sql.Tx{})

			// Check if the response is equal to the expected value
			assert.Equal(t, response, tt.res)

		})
	}
}
