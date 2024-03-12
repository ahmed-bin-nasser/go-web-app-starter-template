package main

// context keys
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

const (
	authenticatedUserIDKey = "authenticatedUserID"
	timestampKey           = "timestamp"
	sessionsKey            = "sessions"
)
