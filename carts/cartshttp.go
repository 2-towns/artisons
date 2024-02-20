package carts

import (
	"artisons/addresses"
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/orders"
	"artisons/shops"
	"artisons/stats"
	"artisons/tags/tree"
	"artisons/templates"
	"artisons/users"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

func getID(r *http.Request, w http.ResponseWriter) int {
	ctx := users.Context(r, w)

	u, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		return u.ID
	}

	id, err := r.Cookie(cookies.CartID)
	if err == nil {
		cid, err := strconv.ParseInt(id.Value, 10, 64)
		if err == nil {
			return int(cid)
		}
	}

	return 0
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	cid := getID(r, w)
	c := Cart{}

	if cid != 0 {
		var err error
		c, err = Get(ctx, cid)
		if err != nil {
			httperrors.Page(w, r.Context(), err.Error(), 500)
			return
		}
	}

	data := struct {
		Lang  language.Tag
		Shop  shops.Settings
		Tags  []tree.Leaf
		Cart  Cart
		Empty bool
	}{
		lang,
		shops.Data,
		tree.Tree,
		c,
		len(c.Products) == 0,
	}

	var t *template.Template
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		t = templates.Pages["hx-cart"]
	} else {
		t = templates.Pages["cart"]
	}

	w.Header().Set("Content-Type", "text/html")

	if cid > 0 {
		c := httphelpers.NewCookie(cookies.CartID, fmt.Sprintf("%d", cid), int(conf.Cookie.MaxAge))
		http.SetCookie(w, &c)
		r.AddCookie(&c)
	}

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pid := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	qty, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "input:quantity")
		return
	}

	cid := getID(r, w)
	if cid == 0 {
		cid, err = NewCartID(ctx)
		if err != nil {
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}
	}

	err = Add(ctx, cid, pid, int(qty))
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	c := httphelpers.NewCookie(cookies.CartID, fmt.Sprintf("%d", cid), int(conf.Cookie.MaxAge))
	http.SetCookie(w, &c)
	r.AddCookie(&c)

	if r.Header.Get("HX-Current-Url") == "/cart" {
		Handler(w, r)
		return
	}

	if shops.Data.Redirect {
		w.Header().Set("HX-Redirect", "/cart")
		w.Write([]byte(""))
		return
	}

	w.Write([]byte(""))
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pid := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	qty, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "input:quantity")
		return
	}

	cid := getID(r, w)

	if !Exists(ctx, cid) {
		httperrors.HXCatch(w, ctx, "you are not authorized to process this request")
		return
	}

	err = Delete(ctx, cid, pid, int(qty))
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	// if r.Header.Get("HX-Current-Url") == "/cart" {
	// 	Handler(w, r)
	// 	return
	// }

	// w.Write([]byte(""))
	Handler(w, r)
}

func DeliveryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cid := getID(r, w)

	c, err := Get(ctx, cid)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 401)
		return
	}

	if len(c.Products) == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "the cart is empty")
		httperrors.Catch(w, ctx, "you are not authorized to process this request", 401)
		return
	}

	del, err := shops.Deliveries(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := struct {
		Lang       language.Tag
		Shop       shops.Settings
		Tags       []tree.Leaf
		Deliveries []string
	}{
		lang,
		shops.Data,
		tree.Tree,
		del,
	}

	w.Header().Set("Content-Type", "text/html")

	coo := httphelpers.NewCookie(cookies.CartID, fmt.Sprintf("%d", cid), int(conf.Cookie.MaxAge))
	http.SetCookie(w, &coo)
	r.AddCookie(&coo)

	if err := templates.Pages["delivery"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func DeliverySetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cid := getID(r, w)

	c, err := Get(ctx, cid)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 401)
		return
	}

	if len(c.Products) == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "the cart is empty")
		httperrors.HXCatch(w, ctx, "you are not authorized to process this request")
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	del := r.FormValue("delivery")
	b := shops.IsValidDelivery(ctx, del)
	if !b {
		httperrors.HXCatch(w, ctx, "the delivery is invalid")
		return
	}

	err = c.UpdateDelivery(ctx, del)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	coo := httphelpers.NewCookie(cookies.CartID, fmt.Sprintf("%d", cid), int(conf.Cookie.MaxAge))
	http.SetCookie(w, &coo)
	r.AddCookie(&coo)

	w.Header().Set("HX-Redirect", "/cart/address")
	w.Write([]byte(""))
}

func AddressFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	cid := getID(r, w)

	c, err := Get(ctx, cid)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 401)
		return
	}

	if len(c.Products) == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "the cart is empty")
		httperrors.Catch(w, ctx, "you are not authorized to process this request", 401)
		return
	}

	data := struct {
		Lang    language.Tag
		Shop    shops.Settings
		Tags    []tree.Leaf
		Address addresses.Address
		URL     string
	}{
		lang,
		shops.Data,
		tree.Tree,
		c.Address,
		"/cart/address",
	}

	coo := httphelpers.NewCookie(cookies.CartID, fmt.Sprintf("%d", cid), int(conf.Cookie.MaxAge))
	http.SetCookie(w, &coo)
	r.AddCookie(&coo)

	if err := templates.Pages["address"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AddressHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cid := getID(r, w)
	log.Println("cid", cid)
	c, err := Get(ctx, cid)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 401)
		return
	}

	if len(c.Products) == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "the cart is empty")
		httperrors.Catch(w, ctx, "you are not authorized to process this request", 401)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	a := addresses.Address{
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

	if err := c.SaveAddress(ctx, a); err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	coo := httphelpers.NewCookie(cookies.CartID, fmt.Sprintf("%d", cid), int(conf.Cookie.MaxAge))
	http.SetCookie(w, &coo)
	r.AddCookie(&coo)

	w.Header().Add("HX-Redirect", "/payment")
}

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	cid := getID(r, w)

	c, err := Get(ctx, cid)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 400)
		return
	}

	err = c.Validate(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 400)
		return
	}

	pay, err := shops.Payments(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	total, err := c.CalculateTotal(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	data := struct {
		Lang     language.Tag
		Shop     shops.Settings
		Tags     []tree.Leaf
		Cart     Cart
		Payments []string
		Total    float64
	}{
		lang,
		shops.Data,
		tree.Tree,
		c,
		pay,
		total,
	}

	coo := httphelpers.NewCookie(cookies.CartID, fmt.Sprintf("%d", cid), int(conf.Cookie.MaxAge))
	http.SetCookie(w, &coo)
	r.AddCookie(&coo)

	if err := templates.Pages["payment"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func PaymentProcessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	payment := r.FormValue("payment")
	cid := getID(r, w)

	c, err := Get(ctx, cid)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 400)
		return
	}

	c.Payment = payment
	log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!", payment)

	err = c.Validate(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 400)
		return
	}

	o := orders.Order{
		Delivery:     c.Delivery,
		DeliveryFees: c.DeliveryFees,
		Payment:      c.Payment,
		Address:      c.Address,
		Products:     c.Products,
		Total:        c.Total,
	}

	err = o.AssignID(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 400)
		return
	}

	redirect, err := shops.Pay(ctx, o.ID, payment)
	if err != nil {
		throws := ctx.Value(contexts.ThrowsWhenPaymentFailed).(bool)
		if throws {
			slog.LogAttrs(ctx, slog.LevelInfo, "the config does not allow to continue when the payment fail")
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}

		o.Status = "payment_progress"
		slog.LogAttrs(ctx, slog.LevelError, "the payment did not work so the payment status is payment_validated")
	} else {
		o.Status = "payment_validated"
	}

	err = o.Save(ctx, c.ID)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 400)
		return
	}

	go o.SendConfirmationEmail(ctx)

	go stats.Order(ctx, o.ID, o.Products, o.Total)

	if redirect != "" {
		w.Header().Add("HX-Redirect", redirect)
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Lang           language.Tag
		SuccessMessage string
	}{
		lang,
		"the order is created successfully",
	}

	if err := templates.Pages["hx-success"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
