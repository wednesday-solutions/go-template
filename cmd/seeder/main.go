package main

import (
	"fmt"
	"go-template/internal/config"
	"go-template/pkg/utl/zaplog"
	"log"
	"os"
	"os/exec"

	"golang.org/x/exp/slices"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		fmt.Println("erorr while loading the env")
		return
	}
	base, _ := os.Getwd()
	base += "/cmd/seeder"
	fmt.Println(base)
	files, err := os.ReadDir(base)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if slices.Contains([]string{"main.go", "utls"}, file.Name()) {
			continue
		}

		filepath := base + "/" + file.Name()
		files, err := os.ReadDir(filepath)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			fmt.Println(filepath + "/" + file.Name())
			cmd := exec.Command("go", "run", filepath+"/"+file.Name())
			data, err := cmd.CombinedOutput()
			if err != nil {
				zaplog.Logger.Error(string(data), err)
			}
			fmt.Println(err)
		}

	}
}
