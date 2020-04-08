package main

import (
	"fmt"
	goboiler "github.com/wednesday-solution/go-boiler"
	"github.com/wednesday-solution/go-boiler/pkg/utl/postgres"
	"github.com/wednesday-solution/go-boiler/pkg/utl/secure"
)

func main() {

	sec := secure.New(1, nil)
	db := postgres.Connect()
	// getting the latest location company and role id so that we can seed a new user
	var location goboiler.Location
	var company goboiler.Company
	var role goboiler.Role
	db.QueryOne(&location, "select id from locations order by id DESC limit 1;")
	db.QueryOne(&company, "select id from companies order by id DESC limit 1;")
	db.QueryOne(&role, "select id from roles order by id ASC limit 1;")
	var insertQuery = fmt.Sprintf("INSERT INTO public.users (first_name, last_name, username, password, email, active, role_id, company_id, location_id) VALUES ('Admin', 'Admin', 'admin', '%s', 'johndoe@mail.com', true, %d, %d, %d);", sec.Hash("adminuser"), role.ID, company.ID, location.ID)
	SeedData("users", insertQuery)
}