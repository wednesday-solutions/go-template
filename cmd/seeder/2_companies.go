package main

func main() {
	SeedData("companies", `INSERT INTO public.companies (name, active) VALUES ('admin_company', true);`)
}
