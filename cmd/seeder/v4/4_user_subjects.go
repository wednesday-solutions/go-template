package main

import (
	"fmt"

	seeders "github.com/wednesday-solutions/go-template/cmd/seeder"
)

func main() {
	var insertQuery = fmt.Sprintf("INSERT INTO user_subjects(user_id, subject_id)" +
		"values(1,1), (1,2), (1,3);")
	_ = seeders.SeedData("user_subjects", insertQuery)
}
