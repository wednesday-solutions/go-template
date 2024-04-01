package config_test

import (
	"os"
	"testing"

	"go-template/internal/config"
	"go-template/testutls"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

const SuccessCase = "Success"

func getLoadTestCases() []struct {
	name     string
	wantData *config.Configuration
	wantErr  bool
	errKey   string
	error    string
} {
	cases := []struct {
		name     string
		wantData *config.Configuration
		wantErr  bool
		errKey   string
		error    string
	}{
		{
			name:     SuccessCase,
			wantErr:  false,
			wantData: testutls.MockConfig(),
		},
		{
			name:    "Failure__NO_SERVER_PORT",
			wantErr: true,
			errKey:  "SERVER_PORT",
			error:   "error loading port from .env",
		},
		{
			name:    "Failure__NO_DB_TIMEOUT_SECONDS",
			wantErr: true,
			errKey:  "DB_TIMEOUT_SECONDS",
			error:   "error loading db timeout from .env",
		},
		{
			name:    "Failure__NO_JWT_MIN_SECRET_LENGTH",
			wantErr: true,
			errKey:  "JWT_MIN_SECRET_LENGTH",
			error:   "error loading jwt min secret length from .env",
		},
		{
			name:    "Failure__NO_APP_MIN_PASSWORD_STR",
			wantErr: true,
			errKey:  "APP_MIN_PASSWORD_STR",
			error:   "error loading application password string from .env",
		},
		{
			name:    "Failure__NO_SERVER_READ_TIMEOUT",
			wantErr: true,
			errKey:  "SERVER_READ_TIMEOUT",
			error:   "error loading server timeout from .env",
		},
	}
	return cases
}

func TestLoad(t *testing.T) {
	cases := getLoadTestCases()
	_, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				patches := ApplyFunc(os.Getenv, func(key string) string {
					if key == tt.errKey {
						return ""
					}
					return key
				})
				defer patches.Reset()
			}
			config, err := config.Load()
			if tt.wantData != nil {
				assert.Equal(t, config, tt.wantData)
			}
			isError := err != nil
			assert.Equal(t, tt.wantErr, isError)
			if isError {
				assert.Equal(t, err.Error(), tt.error)
			}
		})
	}
}
