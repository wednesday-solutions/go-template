package postgres_test

import (
	"database/sql"
	"fmt"
	"go-template/internal/postgres"
	"go-template/testutls"
	"os"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

func TestGetDSN(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "get DSN",
			want: "dbname=go_template host=localhost user=go_template_role password=go_template_role456 port=5432 sslmode=disable",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutls.SetupEnv("../../.env.local")
			if got := postgres.GetDSN(); got != tt.want {
				t.Errorf("GetDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnect(t *testing.T) {
	tests := []struct {
		name    string
		want    *sql.DB
		useOtel bool
		wantErr bool
	}{
		{
			name:    "Return nil when hitting sql.Open",
			want:    nil,
			useOtel: false,
			wantErr: false,
		},
		{
			name:    "Return nil when hitting otelsql.Open",
			want:    nil,
			useOtel: true,
			wantErr: false,
		},
		{
			name:    "Return err when hitting otelsql.Open",
			want:    nil,
			useOtel: true,
			wantErr: true,
		},
	}
	mock, db, _ := testutls.SetupEnvAndDB(t, testutls.Parameters{
		EnvFileLocation: "../../.env.local",
	})
	fmt.Println(mock, db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useOtel {
				os.Args = []string{"something"}
				ApplyFunc(otelsql.Open, func(string, string, ...otelsql.Option) (*sql.DB, error) {
					if tt.wantErr {
						return nil, fmt.Errorf("this is an error")
					}
					return tt.want, nil
				})
			} else {
				ApplyFunc(sql.Open, func(string, string) (*sql.DB, error) {
					return tt.want, nil
				})
			}

			got, err := postgres.Connect()
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connect() = %v, want %v", got, tt.want)
			}
		})
	}
}
