// Package locales provides locale resources for languages
package locales

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/validators"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Value struct {
	Locale string `validate:"required,bcp47_language_tag"`
	Key    string `validate:"required"`
	Value  string `validate:"required"`
}

var UILocale map[string]map[string]string

func init() {
	ctx := context.Background()

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, locale := range conf.LocalesSupported {
			key := "locale:" + locale.String()
			rdb.HGetAll(ctx, key)
		}

		return nil
	})

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the locales", slog.String("error", err.Error()))
		log.Panicln((err))
	}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag links", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		l := strings.Replace(key, "locale:", "", 1)
		tag, err := language.Parse(l)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the language tag", slog.String("tag", l), slog.String("error", err.Error()))
			log.Panicln((err))
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		for k, v := range val {
			message.SetString(tag, k, v)
		}
	}
}

var trans = map[language.Tag]*message.Printer{
	language.English: message.NewPrinter(language.English),
}

// Console is the default language for console
var Console language.Tag = language.English

// Middleware load the detected language in the context.
// It looks into Accept-Language header and fallback
// to english language when the detected language is
// missing or not recognized.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		langs := strings.Split(r.Header.Get("Accept-Language"), "-")
		lang := langs[0]
		var tag language.Tag

		if !slices.Contains(conf.Languages, lang) {
			tag = conf.DefaultLocale
		} else {
			tag = language.Make(lang)
		}

		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		ctx := context.WithValue(r.Context(), contexts.Locale, tag)

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Translate(l language.Tag, msg string, attr ...interface{}) string {
	t := trans[l]

	if t == nil {
		t = trans[conf.DefaultLocale]
	}

	if strings.HasPrefix(msg, "dynamic") {
		return t.Sprintf(fmt.Sprintf("%s%s", msg, attr[0]))
	}

	return t.Sprintf(msg, attr...)
}

func (v Value) Validate(c context.Context) error {
	slog.LogAttrs(c, slog.LevelInfo, "validating a translation")

	if err := validators.V.Struct(v); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the translation", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input_%s_invalid", low)
	}

	slog.LogAttrs(c, slog.LevelInfo, "translation validated")

	return nil
}

func (v Value) Save(ctx context.Context) error {
	l := slog.With(slog.String("key", v.Key))
	l.LogAttrs(ctx, slog.LevelInfo, "saving a translation")

	key := "locale:" + v.Locale

	tag, err := language.Parse(v.Locale)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the locale", slog.String("error", err.Error()))
		return errors.New("error_http_general")
	}

	if _, err := db.Redis.HSet(ctx, key,
		db.Escape(v.Key), db.Escape(v.Value),
	).Result(); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot store the translation", slog.String("error", err.Error()))
		return errors.New("error_http_general")
	}

	message.SetString(tag, v.Key, v.Value)

	l.LogAttrs(ctx, slog.LevelInfo, "translation saved and updated")

	return nil
}
