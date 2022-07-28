package config_test

import (
	"fmt"
	"os"
	"testing"

	"go-template/internal/config"
	"go-template/testutls"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cases := []struct {
		name     string
		wantData *config.Configuration
		wantErr  bool
		errKey   string
		error    string
	}{
		{
			name:     "Success",
			wantErr:  false,
			wantData: testutls.MockConfig(),
		},
		{
			name:    "Failure__NO_SERVER_PORT",
			wantErr: true,
			errKey:  "SERVER_PORT",
			error:   "error loading port from .env ",
		},
		{
			name:    "Failure__NO_DB_TIMEOUT_SECONDS",
			wantErr: true,
			errKey:  "DB_TIMEOUT_SECONDS",
			error:   "error loading db timeout from .env ",
		},
		{
			name:    "Failure__NO_JWT_MIN_SECRET_LENGTH",
			wantErr: true,
			errKey:  "JWT_MIN_SECRET_LENGTH",
			error:   "error loading jwt min secret length from .env ",
		},
		{
			name:    "Failure__NO_APP_MIN_PASSWORD_STR",
			wantErr: true,
			errKey:  "APP_MIN_PASSWORD_STR",
			error:   "error loading application password string from .env ",
		},
		{
			name:    "Failure__NO_SERVER_READ_TIMEOUT",
			wantErr: true,
			errKey:  "SERVER_READ_TIMEOUT",
			error:   "error loading server timeout from .env ",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			_, _, err := testutls.SetupEnvAndDB(t, testutls.Parameters{EnvFileLocation: "../../.env.local"})
			if err != nil {
				fmt.Print("error loading .env file")
			}

			if tt.wantErr == true {
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
