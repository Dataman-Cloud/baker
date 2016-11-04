package cache

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/Dataman-Cloud/baker/model"
)

type Cache struct {
	Users map[string]*model.User
}

func GetUserLogin(c *gin.Context, token string) (*model.User, error) {
	caches := c.MustGet("cache").(Cache)
	for _, v := range caches.Users {
		if v.Token == token {
			return v, nil
		}
	}
	return nil, errors.New("No UserLogin in cache.")
}

func SetUserLogin(c *gin.Context, user *model.User) {
	caches := c.MustGet("cache").(Cache)
	users := caches.Users
	login := user.Login
	users[login] = user
}
