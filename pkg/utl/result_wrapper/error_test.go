package resultwrapper_test

import (
	"github.com/stretchr/testify/assert"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
	"testing"
)

func TestSplitByLabel(t *testing.T) {

	cases := []struct {
		name     string
		req      string
		wantResp string
	}{
		{
			name:     "Error string",
			req:      "no rows in sql",
			wantResp: "no rows in sql",
		},
		{
			name:     "Having `Error` in string",
			req:      `"Error":{"msg"}`,
			wantResp: "\":{\"msg\"}",
		},
		{
			name:     "Having `message` in string",
			req:      `"message":{"msg"}`,
			wantResp: "\":{\"msg\"}",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp := resultwrapper.SplitByLabel(tt.req)
			if len(tt.wantResp) != 0 {
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestErrorFormatter(t *testing.T) {

	cases := []struct {
		name     string
		req      string
		wantResp resultwrapper.ErrorMsg
	}{
		{
			name:     "No string",
			req:      "",
			wantResp: resultwrapper.ErrorMsg{Errors: []string{""}},
		},
		{
			name:     "Having Error string",
			req:      `error message`,
			wantResp: resultwrapper.ErrorMsg{Errors: []string{"error message"}},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp := resultwrapper.ErrorFormatter(tt.req)
			assert.Equal(t, tt.wantResp, resp)
		})
	}
}
