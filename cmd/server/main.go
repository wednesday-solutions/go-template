package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rs/zerolog"

	"go-template/internal/config"

	"go-template/pkg/api"
)

func main() {
	Setup()
}
func Setup() {
	err := config.LoadEnv()
	if err != nil {
		log.Println(err)
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
