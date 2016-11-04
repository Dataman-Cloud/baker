package store

import (
	"golang.org/x/net/context"
)

type Store interface {
	// Authenticate a user by username and password.
	Authenticate(username, password string) (bool, error)
}

// Authenticate a user by username and password.
func Authenticate(c context.Context, username, password string) (bool, error) {
	return FromContext(c).Authenticate(username, password)
}
