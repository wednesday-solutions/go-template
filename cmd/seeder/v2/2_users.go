package main

import (
	"context"
	"fmt"
	"log"

	"go-template/cmd/seeder/utls"
	"go-template/internal/postgres"
	"go-template/models"
	"go-template/pkg/utl/secure"

	"github.com/joho/godotenv"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func main() {

	sec := secure.New(1, nil)
	err := godotenv.Load("./.env.local")
	if err != nil {
		fmt.Print(err)
	}
	db, err := postgres.Connect()
	if err != nil {
		log.Fatal(err, "pop")
	}
	// getting the latest location company and role id so that we can seed a new user

	role, _ := models.Roles(qm.OrderBy("id ASC")).One(context.Background(), db)
	var insertQuery = fmt.Sprintf("INSERT INTO public.users (first_name, last_name, username, password, "+
		"email, active, role_id) VALUES ('Mohammed Ali', 'Chherawalla', 'admin', '%s', 'johndoe@mail.com', true, %d);",
		sec.Hash("adminuser"), role.ID)
	_ = utls.SeedData("users", insertQuery)
}
