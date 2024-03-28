package config_test

import (
	"fmt"
	. "go-template/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		success bool
	}{
		{
			name: "Failed to fetch value from env var",
			args: args{
				key: "key_arg",
			},
			want:    "",
			success: false,
		},
		{
			name: "Successfully getting env var",
			args: args{
				key: "key_arg",
			},
			want:    "value",
			success: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.success {
				os.Setenv(tt.args.key, tt.want)
			}

			if got := GetString(tt.args.key); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		success bool
	}{
		{
			name: "Failed to fetch value from env var",
			args: args{
				key: "int_arg",
			},
			want:    0,
			success: false,
		},
		{
			name: "Successfully getting env var",
			args: args{
				key: "int_arg",
			},
			want:    100,
			success: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.success {
				os.Setenv(tt.args.key, fmt.Sprintf("%d", tt.want))
			}
			if got := GetInt(tt.args.key); got != tt.want {
				t.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		success bool
	}{
		{
			name: "Failed to fetch value from env var",
			args: args{
				key: "bool_arg",
			},
			want:    false,
			success: false,
		},
		{
			name: "Successfully getting env var",
			args: args{
				key: "bool_arg",
			},
			want:    true,
			success: true,
		},
	}
	for _, tt := range tests {
		if tt.success {
			os.Setenv(tt.args.key, fmt.Sprintf("%v", tt.want))
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBool(tt.args.key); got != tt.want {
				t.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Successfully get .env.local",
			want: ".env.local",
		},
		{
			name: "nil env",
			want: ".env",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "nil env" {
				os.Setenv("ENVIRONMENT_NAME", "")
			} else {
				os.Setenv("ENVIRONMENT_NAME", "local")
			}

			if got := FileName(); got != tt.want {
				t.Errorf("FileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

type keyValueArgs struct {
	key   string
	value string
}
type args struct {
	setEnv            []keyValueArgs
	expectedKeyValues []keyValueArgs
}

func TestLoadEnv(t *testing.T) {
	username := "go_template_role"
	host := "localhost"
	dbname := "go_template"
	password := "go_template_role456"
	port := "5432"
	tests := getTestCases(username, host, dbname, password, port)
	for _, tt := range tests {
		setEnvironmentVariables(tt.args)

		t.Run(tt.name, func(t *testing.T) {
			testLoadEnv(t, tt)
		})
	}
}

type envTestCaseArgs struct {
	name    string
	wantErr bool
	args    args
}

func loadLocalEnvIfNoEnvName() envTestCaseArgs {
	return envTestCaseArgs{
		name:    "Successfully load local env if ENVIRONMENT_NAME doesn't have a value",
		wantErr: false,
		args: args{
			setEnv: []keyValueArgs{
				{
					key:   "ENVIRONMENT_NAME",
					value: "",
				},
			},
			expectedKeyValues: []keyValueArgs{
				{
					key:   "PSQL_USER",
					value: "go_template_role",
				},
			},
		},
	}
}

func loadLocalEnv() envTestCaseArgs {
	return envTestCaseArgs{
		name:    "Successfully load local env",
		wantErr: false,
		args: args{
			setEnv: []keyValueArgs{
				{
					key:   "ENVIRONMENT_NAME",
					value: "local",
				},
			},
			expectedKeyValues: []keyValueArgs{
				{
					key:   "SERVER_PORT",
					value: "9000",
				},
			},
		},
	}
}
func errorOnEnvInjectionAndCopilotFalse() envTestCaseArgs {
	return envTestCaseArgs{
		name:    "Error when ENV_INJECTION and COPILOT_DB_CREDS_VIA_SECRETS_MANAGER false",
		wantErr: true,
		args: args{
			setEnv: []keyValueArgs{
				{
					key:   "ENV_INJECTION",
					value: "true",
				},
				{
					key:   "ENVIRONMENT_NAME",
					value: "develop",
				},
				{
					key:   "COPILOT_DB_CREDS_VIA_SECRETS_MANAGER",
					value: "false",
				},
			},
		},
	}
}

func loadOnDbCredsInjected(username string, host string, dbname string, password string, port string) envTestCaseArgs {
	return envTestCaseArgs{
		name:    "dbCredsInjected True",
		wantErr: true,
		args: args{
			setEnv: []keyValueArgs{
				{
					key:   "ENV_INJECTION",
					value: "true",
				},
				{
					key:   "ENVIRONMENT_NAME",
					value: "develop",
				},
				{
					key:   "COPILOT_DB_CREDS_VIA_SECRETS_MANAGER",
					value: "true",
				},
				{
					key: "DB_SECRET",
					value: fmt.Sprintf(`{"username": "%s", "password": "%s", "port": "%s", "dbname": "%s", "host": "%s"}`,
						username,
						password,
						port,
						host,
						dbname),
				},
			},
			expectedKeyValues: []keyValueArgs{
				{
					key:   "PSQL_PASS",
					value: password,
				},
				{
					key:   "PSQL_DBNAME",
					value: dbname,
				},
				{
					key:   "PSQL_HOST",
					value: host,
				},
				{
					key:   "PSQL_USER",
					value: username,
				},
				{
					key:   "PSQL_PORT",
					value: port,
				},
			},
		},
	}
}

func errorOnWrongEnvName() envTestCaseArgs {
	return envTestCaseArgs{
		name:    "Failed to load env",
		wantErr: true,
		args: args{
			setEnv: []keyValueArgs{
				{
					key:   "ENVIRONMENT_NAME",
					value: "local1",
				},
			},
		},
	}
}
func getTestCases(username string, host string, dbname string, password string, port string) []envTestCaseArgs {
	return []envTestCaseArgs{
		loadLocalEnvIfNoEnvName(),
		loadLocalEnv(),
		errorOnEnvInjectionAndCopilotFalse(),
		loadOnDbCredsInjected(username, host, dbname, password, port),
		errorOnWrongEnvName(),
	}
}

func setEnvironmentVariables(args args) {
	for _, env := range args.setEnv {
		os.Setenv(env.key, env.value)
	}
}

func testLoadEnv(t *testing.T, tt struct {
	name    string
	wantErr bool
	args    args
}) {
	if err := LoadEnv(); (err != nil) != tt.wantErr {
		t.Errorf("LoadEnv() error = %v, wantErr %v", err, tt.wantErr)
	} else {
		for _, expected := range tt.args.expectedKeyValues {
			assert.Equal(t, os.Getenv(expected.key), expected.value)
		}
	}
}
