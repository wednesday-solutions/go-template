package utls

import (
	"fmt"
	"go-template/internal/config"
	"go-template/internal/postgres"
	"go-template/pkg/utl/zaplog"
	"log"
	"strings"
)

// SeedData ...
func SeedData(tableName string, rawQuery string) error {
	err := config.LoadEnv()
	if err != nil {
		fmt.Print(err)
	}
	db, err := postgres.Connect()

	if err != nil {
		zaplog.Logger.Error("error while connecting to seed data", err)
		return err
	}
	fmt.Printf("\n-------------------------------\n***Seeding %s\n", tableName)

	queries := strings.Split(rawQuery, ";")

	for _, v := range queries[0 : len(queries)-1] {
		_, err := db.Exec(v)
		if err != nil {
			zaplog.Logger.Error("error while executing seed script for", tableName, err)
			log.Fatal(err)

			return err
		}
	}
	fmt.Printf("***Done seeding %s\n-------------------------------\n", tableName)
	return nil
}
