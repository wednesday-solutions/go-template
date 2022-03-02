package main

import (
	"context"
	"fmt"

	seeders "go-template/cmd/seeder"
	"go-template/internal/postgres"
	"go-template/models"
	"go-template/pkg/utl/secure"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func main() {

	sec := secure.New(1, nil)
	db, _ := postgres.Connect()
	// getting the latest location company and role id so that we can seed a new user

	role, _ := models.Roles(qm.OrderBy("id DESC")).One(context.Background(), db)
	var insertQuery = fmt.Sprintf("INSERT INTO public.users (first_name, last_name, username, password, "+
		"email, active, role_id) VALUES ('Admin', 'Admin', 'admin', '%s', 'johndoe@mail.com', true, %d);",
		sec.Hash("adminuser"), role.ID)
	_ = seeders.SeedData("users", insertQuery)
}
