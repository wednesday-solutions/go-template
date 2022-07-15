package main

import (
	"context"
	"fmt"

	seeders "go-template/cmd/seeder"
	"go-template/internal/mysql"
	"go-template/models"
	"go-template/pkg/utl/secure"
)

func main() {

	sec := secure.New(1, nil)
	db, _ := ms.Connect()
	// getting the latest location company and role id so that we can seed a new user

	_, err := models.Roles().One(context.Background(), db)
	fmt.Println("role is ", err)
	var insertQuery = fmt.Sprintf("INSERT INTO users (first_name, last_name, username, password, "+
		"email, active, role_id) VALUES ('Admin', 'Admin', 'admin', '%s', 'johndoe@mail.com', true, %d);",
		sec.Hash("adminuser"), 1)
	_ = seeders.SeedData("users", insertQuery)
}
