package model

// User represents a user.
type User struct {
	// Login is the username for this user.
	Login string

	// Token is the jwt token.
	Token string

	// Expiry is the token and secret expriation timestamp.
	Expiry int64

	// Hash is a unique token used to sign tokens.
	Hash string

	// Admin user.
	Admin bool
}
