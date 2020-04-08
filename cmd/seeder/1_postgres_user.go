package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"ptcore/pkg/utl/secure"

	"github.com/go-pg/pg"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func random(min int, max int) int {
	return min + rand.Intn(max-min)
}

func main() {

	var psn = `postgres://paytabs_core_dev_role:paytabscore123@localhost:5432/paytabs_core_dev?sslmode=disable`

	u, err := pg.ParseURL(psn)
	checkErr(err)
	db := pg.Connect(u)


	// fmt.Println(a)
	stamp := fmt.Sprint(time.Now().Format("2006-01-02"))
	sec := secure.New(1, nil)
	var password = "admin"
	fmt.Println(password)
	encrypted := sec.Hash(password)
	query := fmt.Sprint("INSERT INTO public.users (mobile_number, created_at, first_name, last_name, email,  age, nationality, password, passcode) VALUES (", random(5000000000, 9999999999), ",", "'", stamp, "'", ",", "'", randSeq(7), "'", ",", "'", randSeq(5), "'", ",", "'", "admin@mail.com", "'", ",", random(2, 150), ",", "'", randSeq(5), "'", ",", "'", encrypted, "'", ",", "'", randSeq(5), "'", ")")
	fmt.Println(query)

	_, err = db.Exec(query)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}
