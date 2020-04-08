package main

import (
	"fmt"
	"github.com/wednesday-solution/go-boiler/pkg/utl/postgres"
	"log"
	"strings"
)

func SeedData(tableName string, rawQuery string) error {
	db := postgres.Connect()

	fmt.Print(fmt.Sprintf("\n-------------------------------\n***Seeding %s\n", tableName))

	queries := strings.Split(rawQuery, ";")

	for _, v := range queries[0 : len(queries)-1] {
		_, err := db.Exec(v)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	fmt.Print(fmt.Sprintf("***Done seeding %s\n-------------------------------\n", tableName))
	return nil
}