package api

import (
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/Dataman-Cloud/baker/cache"
	"github.com/Dataman-Cloud/baker/model"
	"github.com/Dataman-Cloud/baker/store"
)

// Login authenticates the session and returns the
// remote user details.
func Login(c *gin.Context, username, password string) (*model.User, error) {
	// if the username or password is empty, return error.
	if len(username) == 0 || len(password) == 0 {
		logrus.Errorf("username or password is empty.")
		c.AbortWithError(http.StatusUnauthorized, errors.New("username or password is empty"))
		return nil, errors.New("username or password is empty.")
	}
	staticUsersStore := store.FromContext(c)
	if _, err := staticUsersStore.Authenticate(username, password); err != nil {
		logrus.Errorf("authorize failed.")
		c.AbortWithError(http.StatusUnauthorized, errors.New("authenticate failed."))
		return nil, errors.New("authorize failed.")
	}
	return &model.User{Login: username}, nil
}

// Auth authenticates the session and returns the remote user
// login for the given token and secret
func Auth(c *gin.Context, token, secret string) (*model.User, error) {
	cache := c.MustGet("cache").(cache.Cache)
	user, err := cache.GetUserLogin(token)
	if err != nil {
		logrus.Errorf("authorize failed.")
		c.AbortWithError(http.StatusUnauthorized, errors.New("authenticate failed."))
		return nil, errors.New("authorize failed.")
	}
	return user, nil
}
