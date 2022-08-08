package postgres

import (
	"database/sql"
	"fmt"
	"go-template/pkg/utl/zaplog"
	"go-template/testutls"
	"os"

	otelsql "github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	// DB adapter
	_ "github.com/lib/pq"
)

// Connect ...
func Connect() (*sql.DB, error) {
	dsn := GetDSN()
	zaplog.Logger.Info("Connecting to DB\n", dsn)
	if testutls.IsInTests() {
		return sql.Open("postgres", dsn)
	}
	return otelsql.Open("postgres", dsn, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
}

func GetDSN() string {
	dsn := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=%s",
		os.Getenv("PSQL_DBNAME"),
		os.Getenv("PSQL_HOST"),
		os.Getenv("PSQL_USER"),
		os.Getenv("PSQL_PASS"),
		os.Getenv("PSQL_PORT"),
		os.Getenv("PSQL_SSLMODE"))
	return dsn
}
