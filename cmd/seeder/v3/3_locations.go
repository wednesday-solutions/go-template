package main

import seeders "github.com/wednesday-solutions/go-template/cmd/seeder"

func main() {
	_ = seeders.SeedData("locations", `INSERT INTO public.locations (name, active, company_id, address) VALUES ('admin_location', true, 1, 'admin_address');`)
}
