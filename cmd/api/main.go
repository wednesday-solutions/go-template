package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/wednesday-solutions/go-boiler/pkg/api"

	"github.com/wednesday-solutions/go-boiler/pkg/utl/config"
)

func main() {

	err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	cfg, err := config.Load()
	checkErr(err)

	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
