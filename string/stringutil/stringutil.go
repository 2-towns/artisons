// Package stringutil provides string utilities
package stringutil

import (
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Random provides a random unique string
func Random() string {
	length := 16
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

// Slugify returns the slug representation of a title
func Slugify(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}
