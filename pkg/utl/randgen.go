package utl

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var Intn = rand.Intn

// RandomSequence ...
func RandomSequence(n int) string {
	b := make([]rune, n)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = letters[Intn(len(letters))]
	}
	return string(b)
}
