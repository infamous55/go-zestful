package main

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(length uint) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	buffer := make([]rune, length)
	for i := range buffer {
		buffer[i] = letters[random.Intn(len(letters))]
	}
	return string(buffer)
}
