package main

import (
	"fmt"
	"go-template/internal/config"
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
	base := "./cmd/seeder"
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
			err := exec.Command("go", "run", filepath+"/"+file.Name()).Run()
			fmt.Println(err)
		}

	}
}
