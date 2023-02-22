package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/exp/slices"
)

func main() {

	base, _ := os.Getwd()
	log.Println(base)

	base += "/cmd/seeder/exec/build"
	fmt.Println(base)

	files, err := os.ReadDir(base)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(files))

	for _, file := range files {

		log.Println(file.Name())
		if slices.Contains([]string{"seed.go"}, file.Name()) && strings.Contains(file.Name(), "env") {
			continue
		}
		cmd := exec.
			Command(fmt.Sprintf("%s/%s", base, file.Name()))
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("out:", outb.String(), "err:", errb.String())

	}
}
