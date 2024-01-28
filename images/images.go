package images

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"artisons/conf"
	"log/slog"
)

type Options struct {
	Width       string
	Height      string
	Cachebuster int64
}

// URL returns the imgproxy URL
func URL(id string, o Options) string {
	if id == "" {
		return ""
	}

	path := fmt.Sprintf(`/resize:fill:%s:%s/cachebuster:%d/plain/%s/%s`, o.Width, o.Height, o.Cachebuster, conf.ImgProxy.Protocol, id)

	if conf.ImgProxy.Key != "" && conf.ImgProxy.Salt != "" {
		hmac := hmac.New(sha256.New, []byte(conf.ImgProxy.Key))

		sal, err := hex.DecodeString(conf.ImgProxy.Salt)
		if err != nil {
			slog.Error("cannot get the bytes salt value", slog.String("error", err.Error()))
			return ""
		}

		hmac.Write(sal)
		hmac.Write([]byte(path))
		sum := hmac.Sum(nil)
		sig := hex.EncodeToString(sum)

		return fmt.Sprintf("%s%s%s", conf.ImgProxy.URL, sig, path)
	}

	return fmt.Sprintf("%s%s", conf.ImgProxy.URL, path)
}
