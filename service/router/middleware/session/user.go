package session

import (
	_ "net/http"

	_ "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/Dataman-Cloud/baker/cache"
	"github.com/Dataman-Cloud/baker/model"
	"github.com/Dataman-Cloud/baker/service/shared/token"
)

func User(c *gin.Context) *model.User {
	v, ok := c.Get("user")
	if !ok {
		return nil
	}
	u, ok := v.(*model.User)
	if !ok {
		return nil
	}
	return u
}

func Token(c *gin.Context) *token.Token {
	v, ok := c.Get("token")
	if !ok {
		return nil
	}
	u, ok := v.(*token.Token)
	if !ok {
		return nil
	}
	return u
}

func SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *model.User

		_, err := token.ParseRequest(c.Request, func(t *token.Token) (string, error) {
			var err error
			cache := c.MustGet("cache").(cache.Cache)
			user, err = cache.GetUserLogin(t.Text)
			return user.Hash, err
		})
		if err == nil {
			if user.Login == "admin" {
				user.Admin = true
			}
			c.Set("user", user)

			// if this is a session token (ie not the API token)
			// this means the user is accessing with a web browser,
			// so we should implement CSRF protection measures.
			//if t.Kind == token.SessToken {
			//	err = token.CheckCsrf(c.Request, func(t *token.Token) (string, error) {
			//		return user.Hash, nil
			//	})
			//	// if csrf token validation fails, exit immediately
			//	// with a not authorized error.
			//	if err != nil {
			//		c.AbortWithStatus(http.StatusUnauthorized)
			//		return
			//	}
			//}
		}
		c.Next()
	}
}

func MustAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)
		switch {
		case user == nil:
			c.String(401, "User not authorized")
			c.Abort()
		case user.Admin == false:
			c.String(403, "User not authorized")
			c.Abort()
		default:
			c.Next()
		}
	}
}

func MustUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)
		switch {
		case user == nil:
			c.String(401, "User not authorized")
			c.Abort()
		default:
			c.Next()
		}
	}
}
