package main

import seeders "github.com/wednesday-solutions/go-template/cmd/seeder"

func main() {
	_ = seeders.SeedData("companies", `INSERT INTO public.companies (name, active) VALUES ('admin_company', true);`)
}
