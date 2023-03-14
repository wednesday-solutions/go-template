package config

import (
	"encoding/json"
	"fmt"
	"go-template/pkg/utl/convert"
	"log"
	"os"
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

func LoadEnvWithFilePrefix(fileprefix *string) error {
	prefix := ""
	if fileprefix != nil {
		prefix = *fileprefix
	}
	err := godotenv.Load(fmt.Sprintf("%s.env.base", prefix))
	if err != nil {
		return err
	}

	envName := os.Getenv("ENVIRONMENT_NAME")
	if envName == "" {
		envName = "local"
	}
	log.Println("envName: " + envName)

	envVarInjection := GetBool("ENV_INJECTION")
	if !envVarInjection || envName == "local" {
		err = godotenv.Load(fmt.Sprintf("%s.env.%s", prefix, envName))

		if err != nil {
			fmt.Printf(".env.%s\n", envName)
			log.Println(err)
			return err
		}
	}

	dbCredsInjected := GetBool("COPILOT_DB_CREDS_VIA_SECRETS_MANAGER")

	if dbCredsInjected {
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
			return err
		}

		os.Setenv("PSQL_DBNAME", secrets.DBName)
		os.Setenv("PSQL_HOST", secrets.Host)
		os.Setenv("PSQL_PASS", secrets.Password)
		os.Setenv("PSQL_PORT", strconv.Itoa(secrets.Port))
		os.Setenv("PSQL_USER", secrets.Username)

	}

	return nil
}
func LoadEnv() error {
	return LoadEnvWithFilePrefix(nil)
}
