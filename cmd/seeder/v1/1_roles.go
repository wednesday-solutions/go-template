package main

import "go-template/cmd/seeder/utls"

func main() {
	_ = utls.SeedData("roles", `INSERT INTO public.roles("access_level", "name") VALUES (100, 'SUPER_ADMIN');
		INSERT INTO public.roles("access_level", "name") VALUES (120, 'COMPANY_ADMIN');
		INSERT INTO public.roles("access_level", "name") VALUES (130, 'LOCATION_ADMIN');
		INSERT INTO public.roles("access_level", "name") VALUES (200, 'USER');
		INSERT INTO public.roles("access_level", "name") VALUES (110, 'ADMIN');`)
}
