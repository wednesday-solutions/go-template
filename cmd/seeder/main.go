package seeders

import (
	"fmt"
	"log"
	"strings"

	"go-template/internal/postgres"
)

// SeedData ...
func SeedData(tableName string, rawQuery string) error {
	db, _ := postgres.Connect()

	fmt.Printf("\n-------------------------------\n***Seeding %s\n", tableName)

	queries := strings.Split(rawQuery, ";")

	for _, v := range queries[0 : len(queries)-1] {
		_, err := db.Exec(v)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	fmt.Printf("***Done seeding %s\n-------------------------------\n", tableName)
	return nil
}
