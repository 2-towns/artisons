// Package contexts provide utilities related to HTTP
package contexts

// ContextKey is the type of key used for the context.
// It is necessary to create a specific type for the context, but
// it does not bring added value.
type contextKey string

// User is the context key used to store the lang
const User contextKey = "user"

// User is the context key used to store the lang
const UserID contextKey = "user_id"

// Cart is the context key used to store the lang
const Cart contextKey = "cart"

// End can "front" or "back"
const End contextKey = "end"

// Demo is the context key used when the admin activated the demo mode
const Demo contextKey = "demo"

// Locale is the context key used to store the lang
const Locale contextKey = "lang"

// HX true if the request is htmx request
const HX contextKey = "hx"

// HXTarget is used to change the default target of alert message
const HXTarget contextKey = "hx-target"

const AlertTarget contextKey = "alert"
