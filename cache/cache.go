package cache

import (
	"errors"

	"github.com/Dataman-Cloud/baker/model"
)

type Cache struct {
	Users map[string]*model.User
}

func NewUserLoginCache() Cache {
	return Cache{Users: make(map[string]*model.User)}
}

func (cache *Cache) GetUserLogin(login string) (*model.User, error) {
	for k, v := range cache.Users {
		if k == login {
			return v, nil
		}
	}
	return nil, errors.New("No UserLogin in cache.")
}

func (cache *Cache) SetUserLogin(user *model.User) {
	users := cache.Users
	users[user.Login] = user
}
