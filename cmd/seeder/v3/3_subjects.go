package main

import (
	"fmt"

	seeders "github.com/wednesday-solutions/go-template/cmd/seeder"
)

func main() {
	var insertQuery = fmt.Sprintf("INSERT INTO subjects(name)" +
		"values('ENGLISH'), ('HISTORY'), ('MATHS');")
	_ = seeders.SeedData("subjects", insertQuery)
}
