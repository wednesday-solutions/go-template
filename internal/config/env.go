package config

import (
	"encoding/json"
	"fmt"
	"go-template/pkg/utl/convert"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

func GetString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		keyNotFound(key)
	}

	return value
}

func GetInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		keyNotFound(key)
		return 0
	}
	return convert.StringToInt(value)
}

func GetBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		keyNotFound(key)
		return false
	}
	return convert.StringToBool(value)
}

func keyNotFound(key string) {
	fmt.Printf("Key %s not found in %s Returning default value.", key, FileName())
}

func FileName() string {
	environment := os.Getenv("ENVIRONMENT_NAME")
	var envFileName string

	if environment == "" {
		envFileName = ".env"
	} else {
		envFileName = fmt.Sprintf(".env.%s", environment)
	}
	return envFileName
}

func LoadEnv() error {
	const (
		localEnvFile = "local"
	)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("Error getting current file path")
	}

	prefix := fmt.Sprintf("%s/", filepath.Join(filepath.Dir(filename), "../../"))
	err := godotenv.Load(fmt.Sprintf("%s.env.base", prefix))
	if err != nil {
		return err
	}
	fmt.Println("loaded", fmt.Sprintf("%s.env.base", prefix))

	envName := os.Getenv("ENVIRONMENT_NAME")
	if envName == "" {
		envName = localEnvFile
	}
	log.Println("envName: " + envName)

	envVarInjection := GetBool("ENV_INJECTION")
	if !envVarInjection || envName == localEnvFile {
		err = godotenv.Load(fmt.Sprintf("%s.env.%s", prefix, envName))

		if err != nil {
			return fmt.Errorf("couldn't load .env.%s file: %w", envName, err)
		}
		fmt.Println("loaded", fmt.Sprintf("%s.env.%s", prefix, envName))
		return nil
	}

	dbCredsInjected := GetBool("COPILOT_DB_CREDS_VIA_SECRETS_MANAGER")

	// except for local environment the db creds should be
	// injected through the secret manager
	if envName != localEnvFile && !dbCredsInjected {
		return fmt.Errorf("db creds should be injected through secret manager")
	}

	// for local environment the db configs will be loaded through
	// env itself by `godotenv`, even if db creds are not injected
	if !dbCredsInjected {
		return nil
	}

	// if dbCreds is injected then extract the db creds
	return extractDBCredsFromSecret()
}

// extractDBCredsFromSecret helper function to extract db secret
func extractDBCredsFromSecret() error {
	type copilotSecrets struct {
		Username string `json:"username"`
		Host     string `json:"host"`
		DBName   string `json:"dbname"`
		Password string `json:"password"`
		Port     int    `json:"port"`
	}
	secrets := &copilotSecrets{}

	err := json.Unmarshal([]byte(os.Getenv("DB_SECRET")), secrets)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal db secret: %w", err)
	}

	os.Setenv("PSQL_DBNAME", secrets.DBName)
	os.Setenv("PSQL_HOST", secrets.Host)
	os.Setenv("PSQL_PASS", secrets.Password)
	os.Setenv("PSQL_PORT", strconv.Itoa(secrets.Port))
	os.Setenv("PSQL_USER", secrets.Username)

	return nil
}
