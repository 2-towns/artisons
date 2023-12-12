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

// Demo is the context key used when the admin activated the demo mode
const Demo contextKey = "demo"

// Locale is the context key used to store the lang
const Locale contextKey = "lang"

const AlertTarget contextKey = "alert"
