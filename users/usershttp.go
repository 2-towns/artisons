package users

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/shops"
	"artisons/tags/tree"
	"artisons/templates"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

func findBySessionID(ctx context.Context, sid string) (User, error) {
	l := slog.With(slog.String("sid", sid))
	l.LogAttrs(ctx, slog.LevelInfo, "finding the user")

	if sid == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the session id")
		return User{}, errors.New("you are not authorized to process this request")
	}

	id, err := db.Redis.HGet(ctx, "session:"+sid, "uid").Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the auth id from redis", slog.String("error", err.Error()))
		return User{}, errors.New("you are not authorized to process this request")
	}

	m, err := db.Redis.HGetAll(ctx, "user:"+id).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the session from redis", slog.String("error", err.Error()))
		return User{}, errors.New("you are not authorized to process this request")
	}

	m["sid"] = sid
	u, err := parse(ctx, m)
	if err != nil {
		return User{}, errors.New("you are not authorized to process this request")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "user found", slog.Int("user_id", u.ID))

	return u, err
}

func Context(r *http.Request, w http.ResponseWriter) context.Context {
	ctx := r.Context()

	sid, err := r.Cookie(cookies.SessionID)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "no session cookie found")
		return ctx
	}

	user, err := findBySessionID(ctx, sid.Value)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "session id not found so destroying it", slog.String("sid", sid.Value))

		c := httphelpers.NewCookie(cookies.SessionID, sid.Value, -1)
		http.SetCookie(w, &c)

		return ctx
	} else {
		// Refresh the cookie
		c := httphelpers.NewCookie(cookies.SessionID, sid.Value, int(conf.Cookie.MaxAge))
		http.SetCookie(w, &c)
	}

	ctx = context.WithValue(ctx, contexts.User, user)

	return ctx
}

func AccountOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := Context(r, w)
		user, ok := ctx.Value(contexts.User).(User)

		if !ok {
			slog.LogAttrs(ctx, slog.LevelInfo, "no session cookie found")
			//httperrors.Catch(w, ctx, "you are not authorized to process this request", 401)
			http.Redirect(w, r, "/otp", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))

		err := user.RefreshSession(ctx)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot refresh the session", slog.String("error", err.Error()))
		}
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := Context(r, w)

		user, ok := ctx.Value(contexts.User).(User)

		if !ok {
			slog.LogAttrs(ctx, slog.LevelInfo, "no session cookie found")
			//httperrors.Catch(w, ctx, "you are not authorized to process this request", 401)
			http.Redirect(w, r, "/sso", http.StatusFound)
			return
		}

		if user.Role != "admin" {
			slog.LogAttrs(ctx, slog.LevelInfo, "the user is not admin", slog.Int("id", user.ID))
			httperrors.Catch(w, ctx, "you are not authorized to process this request", 401)
			return
		}

		w.Header().Set("X-Robots-Tag", "noindex")
		next.ServeHTTP(w, r.WithContext(ctx))

		log.Println("admin donene !!!")

		err := user.RefreshSession(ctx)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot refresh the session", slog.String("error", err.Error()))
		}
	})
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Lang language.Tag
		Shop shops.Settings
		Tags []tree.Leaf
	}{
		lang,
		shops.Data,
		tree.Tree,
	}

	if err := templates.Pages["account"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AddressFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	user := ctx.Value(contexts.User).(User)

	data := struct {
		Lang language.Tag
		Shop shops.Settings
		Tags []tree.Leaf
		User User
	}{
		lang,
		shops.Data,
		tree.Tree,
		user,
	}

	if err := templates.Pages["address"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AddressHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	user := ctx.Value(contexts.User).(User)

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	a := Address{
		Firstname:     r.FormValue("firstname"),
		Lastname:      r.FormValue("lastname"),
		Street:        r.FormValue("street"),
		Complementary: r.FormValue("complementary"),
		City:          r.FormValue("city"),
		Zipcode:       r.FormValue("zipcode"),
		Phone:         r.FormValue("phone"),
	}

	if err := a.Validate(ctx); err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if err := a.Save(ctx, user.ID); err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	data := struct {
		Lang           language.Tag
		SuccessMessage string
	}{
		lang,
		"The address has been saved successfully.",
	}

	if err := templates.Pages["hx-success"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
