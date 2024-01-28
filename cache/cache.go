package cache

import (
	"crypto/md5"
	"encoding/hex"
	"artisons/conf"
	"log"
	"log/slog"
	"os"
	"path"
	"strings"
)

var buster = map[string]string{}

func Buster(name string) string {
	return buster[name]
}

func load(folder string, ext string) {
	files, err := os.ReadDir(conf.WorkingSpace + folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		name := f.Name()

		if !strings.HasSuffix(name, ext) {
			continue
		}

		buf, err := os.ReadFile(path.Join(folder, name))

		if err != nil {
			log.Fatal(err)
		}

		hash := md5.Sum(buf)
		h := hex.EncodeToString(hash[:])
		buster[name] = h
	}

	slog.Info("files loaded", slog.String("folder", folder), slog.Int("length", len(files)))
}

func Busting() {
	load("web/public/js/admin", "js")
	load("web/public/css/admin", "css")
}
