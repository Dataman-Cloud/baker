package model

// User represents a user.
type User struct {
	// Login is the username for this user.
	Login string `json:"login"`

	// Token is the jwt token.
	Token string `json:"-"`

	// Expiry is the token and secret expriation timestamp.
	Expiry int64 `json:"-"`

	// Hash is a unique token used to sign tokens.
	Hash string `json:"-"`
}
