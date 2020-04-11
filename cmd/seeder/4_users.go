package main

import (
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/wednesday-solutions/go-boiler/models"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/postgres"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/secure"
)

func main() {

	sec := secure.New(1, nil)
	db, _ := postgres.Connect()
	// getting the latest location company and role id so that we can seed a new user

	location, _ := models.Locations(qm.OrderBy("id DESC")).One(context.Background(), db)
	company, _ := models.Companies(qm.OrderBy("id DESC")).One(context.Background(), db)
	var insertQuery = fmt.Sprintf("INSERT INTO public.users (first_name, last_name, username, password, "+
		"email, active, company_id, location_id) VALUES ('Admin', 'Admin', 'admin', '%s',"+
		" 'johndoe@mail.com', true, %d, %d);",
		sec.Hash("adminuser"),
		company.ID,
		location.ID)
	SeedData("users", insertQuery)
}
