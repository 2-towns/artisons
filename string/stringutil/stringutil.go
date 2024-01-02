// Package stringutil provides string utilities
package stringutil

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"
)

// Random provides a random unique string
func Random() (string, error) {
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// Slugify returns the slug representation of a title
func Slugify(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}
