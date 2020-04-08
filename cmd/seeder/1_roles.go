package main

func main() {
	SeedData("roles", `INSERT INTO public.roles VALUES (110, 110, 'ADMIN');
		INSERT INTO public.roles VALUES (120, 120, 'COMPANY_ADMIN');
		INSERT INTO public.roles VALUES (130, 130, 'LOCATION_ADMIN');
		INSERT INTO public.roles VALUES (200, 200, 'USER');
		INSERT INTO public.roles VALUES (100, 100, 'SUPER_ADMIN');`)
}
