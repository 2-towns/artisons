// Package contexts provide utilities related to HTTP
package contexts

// ContextKey is the type of key used for the context.
// It is necessary to create a specific type for the context, but
// it does not bring added value.
type ContextKey string

// User is the context key used to store the lang
const User ContextKey = "user"

// Form is the current form object
const Form ContextKey = "form"

// Device is the context key used to store the lang
const Device ContextKey = "device"

// Locale is the context key used to store the lang
const Locale ContextKey = "lang"

// HX true if the request is htmx request
const HX ContextKey = "hx"

const RequestID ContextKey = "request-id"

// HXTarget is used to change the default target of alert message
const HXTarget ContextKey = "hx-target"

const AlertTarget ContextKey = "alert"

const Tracking ContextKey = "tracking"
