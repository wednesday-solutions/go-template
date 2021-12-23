package config_test

import (
	"fmt"
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
			err := godotenv.Load("../../.env.local")
			if err != nil {
				fmt.Print("error loading .env file")
			}
			_, err = config.Load()
			fmt.Print("\n", err, "\n")
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
