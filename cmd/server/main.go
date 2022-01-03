package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/joho/godotenv"
	"github.com/wednesday-solutions/go-template/internal/config"
	"github.com/wednesday-solutions/go-template/pkg/api"
)

func main() {
	envName := os.Getenv("ENVIRONMENT_NAME")

	if envName == "" {
		envName = "local"
	}

	err := godotenv.Load(fmt.Sprintf(".env.%s", envName))
	if err != nil {
		fmt.Print("error loading .env file")
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
