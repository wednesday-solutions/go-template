package main

import (
	"github.com/99designs/gqlgen/cmd"
	_ "github.com/web-ridge/sqlboiler-graphql-schema"

)

func main() {
	cmd.Execute()
}
