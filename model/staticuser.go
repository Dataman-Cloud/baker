package model

// StaticUser type.
type StaticUser struct {
	// Username is the username for this user.
	Username string `json:"username"`

	// the password for this user.
	Password string `json:"password"`
}
