package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/wednesday-solutions/go-template/internal/config"
	"github.com/wednesday-solutions/go-template/pkg/api"
)

func main() {

	err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
	if err != nil {
		fmt.Print("Error loading .env file")
		checkErr(err)
		os.Exit(1)
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
