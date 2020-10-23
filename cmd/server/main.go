package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"

	"github.com/joho/godotenv"
	"github.com/wednesday-solutions/go-template/pkg/api"

	"github.com/wednesday-solutions/go-template/pkg/utl/config"
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
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Panic().Msg(err.Error())
	}
}
