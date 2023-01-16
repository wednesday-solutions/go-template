package main

import (
	"fmt"
	"go-template/internal/config"
	"go-template/pkg/utl/zaplog"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/exp/slices"
)

func stripFileExtension(fileName string) string {
	s := strings.TrimRight(fileName, ".go")
	return s
}

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
		if slices.Contains([]string{"main.go", "seed", "output", "exec", "build", "utls"}, file.Name()) {
			continue
		}

		filepath := base + "/" + file.Name()
		// outputPath := base + "/output"
		files, err := os.ReadDir(filepath)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			fmt.Println(filepath + "/" + file.Name())
			cmd := exec.
				Command("go", "build", "-o",
					fmt.Sprintf("./cmd/seeder/exec/build/%s", stripFileExtension(file.Name())), filepath+"/"+file.Name())
			data, err := cmd.CombinedOutput()
			if err != nil {
				zaplog.Logger.Error(string(data), err)
			}
			fmt.Println(err)
		}

	}
}
