package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/wednesday-solutions/go-template/internal/config"
)

func TestLoad(t *testing.T) {
	cases := []struct {
		name     string
		wantData *config.Configuration
		wantErr  bool
	}{
		{
			name:    "Success",
			wantErr: false,
			wantData: &config.Configuration{
				DB: &config.Database{
					LogQueries: true,
					Timeout:    20,
				},
				Server: &config.Server{
					Port:         ":8080",
					Debug:        true,
					ReadTimeout:  15,
					WriteTimeout: 20,
				},
				JWT: &config.JWT{
					MinSecretLength:  128,
					DurationMinutes:  10,
					RefreshDuration:  10,
					MaxRefresh:       144,
					SigningAlgorithm: "HS384",
				},
				App: &config.Application{
					MinPasswordStr: 3,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf("../../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
			if err != nil {
				fmt.Print("Error loading .env file")
			}
			_, err = config.Load()
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
