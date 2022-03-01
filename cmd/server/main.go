package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"go-template/internal/config"
	"go-template/pkg/api"
)

func main() {
	Setup()
}
func Setup() {
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
	_, err = api.Start(cfg)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Panic().Msg(err.Error())
	}
}
