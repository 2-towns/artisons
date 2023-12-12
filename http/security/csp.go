package security

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

var CSP map[string]string = map[string]string{}

func LoadCsp() {
	entries, err := os.ReadDir("web/views/admin/js")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		name := e.Name()

		tpl, err := template.ParseFiles(
			"web/views/admin/js/" + name,
		)
		if err != nil {
			log.Fatal(err)
		}

		b := bytes.NewBufferString("")
		tpl.Execute(b, "")
		clean := strings.Replace(b.String(), "<script>", "", 1)
		clean = strings.Replace(clean, "</script>", "", 1)

		hasher := sha256.New()
		hasher.Write([]byte(clean))
		key := strings.Replace(e.Name(), ".js.html", "", 1)
		hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		CSP[key] = fmt.Sprintf("'sha256-%s'", hash)
	}
}
