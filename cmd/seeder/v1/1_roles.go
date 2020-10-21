package main

import seeders "github.com/wednesday-solutions/go-template/cmd/seeder"

func main() {
	_ = seeders.SeedData("roles", `INSERT INTO public.roles VALUES (1, 110, 'ADMIN');
		INSERT INTO public.roles VALUES (2, 120, 'COMPANY_ADMIN');
		INSERT INTO public.roles VALUES (3, 130, 'LOCATION_ADMIN');
		INSERT INTO public.roles VALUES (4, 200, 'USER');
		INSERT INTO public.roles VALUES (5, 100, 'SUPER_ADMIN');`)
}
