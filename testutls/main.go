package testutls

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"go-template/internal/config"
	"go-template/models"

	"github.com/DATA-DOG/go-sqlmock"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type key string

var (
	UserKey key = "user"
)

const (
	UserRoleName       = "UserRole"
	SuperAdminRoleName = "SuperAdminRole"
	SuccessCase        = "Success"
	ErrorFindingUser   = "Fail on finding user"
)

var MockIpAddress = "0.0.0.0"
var MockEmail = "mac@wednesday.is"
var MockToken = "token_string"
var MockID = 1
var MockCount = int64(1)
var MockJWTSecret = "1234567890123456789012345678901234567890123456789012345678901234567890"
var MockQuery = `{"query":"query users { users { users { id } } }","variables":{}}"`
var MockWhitelistedQuery = `{"query":"query Schema {  __schema { queryType { kind } } }","variables":{}}"`

func MockUser() *models.User {
	return &models.User{
		ID:                 MockID,
		FirstName:          null.StringFrom("First"),
		LastName:           null.StringFrom("Last"),
		Username:           null.StringFrom("username"),
		Email:              null.StringFrom(MockEmail),
		Mobile:             null.StringFrom("+911234567890"),
		Address:            null.StringFrom("22 Jump Street"),
		Password:           null.StringFrom(`!@#!@12ASa3123`),
		Active:             null.BoolFrom(false),
		LastLogin:          null.NewTime(time.Time{}, false),
		LastPasswordChange: null.NewTime(time.Time{}, false),
		Token:              null.StringFrom("asd"),
		DeletedAt:          null.NewTime(time.Time{}, false),
		UpdatedAt:          null.NewTime(time.Time{}, false),
		RoleID:             null.IntFrom(1),
	}
}
func MockUsers() []*models.User {
	return []*models.User{
		{
			FirstName: null.StringFrom("First"),
			LastName:  null.StringFrom("Last"),
			Username:  null.StringFrom("username"),
			Email:     null.StringFrom(MockEmail),
			Mobile:    null.StringFrom("+911234567890"),
			Address:   null.StringFrom("22 Jump Street"),
		},
	}

}

func MockJwt(role string) *jwt.Token {
	return &jwt.Token{
		Raw: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIi" +
			"wibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		Method: jwt.GetSigningMethod("HS256"),
		Claims: jwt.MapClaims{
			"e":    MockEmail,
			"exp":  "1.641189209e+09",
			"id":   MockID,
			"u":    "admin",
			"sub":  "1234567890",
			"name": "John Doe",
			"iat":  1516239022,
			"role": role,
		},
		Header: map[string]interface{}{
			"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		Valid: true,
	}
}

type Parameters struct {
	EnvFileLocation string `default:"../.env.local"`
}

func SetupEnv(envfile string) {
	err := godotenv.Load(envfile)
	if err != nil {
		fmt.Print("error loading .env file")
	}
}
func SetupEnvAndDB(t *testing.T, parameters Parameters) (mock sqlmock.Sqlmock, db *sql.DB, err error) {
	SetupEnv(parameters.EnvFileLocation)
	db, mock, err = sqlmock.New()
	if err != nil {
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
	}
	boil.SetDB(db)
	return mock, db, nil
}

func SetupMockDB(t *testing.T) (mock sqlmock.Sqlmock, db *sql.DB, err error) {
	db, mock, err = sqlmock.New()
	if err != nil {
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
	}
	boil.SetDB(db)
	return mock, db, nil
}

type QueryData struct {
	Actions    *[]driver.Value
	Query      string
	DbResponse *sqlmock.Rows
}

func MockConfig() *config.Configuration {
	return &config.Configuration{
		DB: &config.Database{
			LogQueries: true,
			Timeout:    5,
		},
		Server: &config.Server{
			Port:         ":9000",
			Debug:        true,
			ReadTimeout:  10,
			WriteTimeout: 5,
		},
		JWT: &config.JWT{
			MinSecretLength:  64,
			DurationMinutes:  1440,
			RefreshDuration:  3499200,
			MaxRefresh:       1440,
			SigningAlgorithm: "HS256",
		},
		App: &config.Application{
			MinPasswordStr: 1,
		},
	}
}
func IsInTests() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.paniconexit0") {
			return true
		}
	}
	return false
}
