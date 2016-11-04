package store

import (
	"errors"

	"github.com/Dataman-Cloud/baker/model"
	"golang.org/x/crypto/bcrypt"
)

type StaticUsersStore struct {
	Users map[string]*model.StaticUser
}

func NewStaticUsersStore(users map[string]*model.StaticUser) Store {
	return &StaticUsersStore{Users: users}
}

func (s *StaticUsersStore) Authenticate(username, password string) (bool, error) {
	user := s.Users[username]
	if user == nil {
		return false, errors.New("NoMatch")
	}
	if user != nil {
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
			return false, nil
		}
	}
	return true, nil
}
