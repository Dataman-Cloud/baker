package api

import (
	"encoding/base32"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"

	"github.com/Dataman-Cloud/baker/cache"
	"github.com/Dataman-Cloud/baker/service/shared/httputil"
	"github.com/Dataman-Cloud/baker/service/shared/token"
)

// ShowLogin is a endpoint that redirects to
// initiliaze the oauth flow
func ShowLogin(c *gin.Context) {
	c.Redirect(303, "/authorize")
}

func GetLogin(c *gin.Context) {
	// when dealing with redirects we may need to adjust the content type. I
	// cannot, however, remember why, so need to revisit this line.
	c.Writer.Header().Del("Content-Type")
	in := &loginPayload{}
	err := c.Bind(in)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := Login(c, in.Username, in.Password)
	if err != nil {
		logrus.Errorf("cannot authenticate user. %s", err)
		c.Redirect(303, "/login?error=oauth_error")
		return
	}

	// construct a login user into cache.
	user.Hash = base32.StdEncoding.EncodeToString(
		securecookie.GenerateRandomKey(32),
	)

	exp := time.Now().Add(time.Hour * 72).Unix()
	token := token.New(token.SessToken, user.Login)
	tokenstr, err := token.SignExpires(user.Hash, exp)
	if err != nil {
		logrus.Errorf("cannot create token for %s. %s", user.Login, err)
		c.Redirect(303, "/login?error=internal_error")
		return
	}
	// update userLogin in Context.
	user.Token = tokenstr
	user.Expiry = exp
	cache := c.MustGet("cache").(cache.Cache)
	cache.SetUserLogin(user)

	httputil.SetCookie(c.Writer, c.Request, "user_sess", tokenstr)
	c.JSON(http.StatusOK, &tokenPayload{
		Access:  tokenstr,
		Expires: exp - time.Now().Unix(),
	})
}

func GetLogout(c *gin.Context) {
	httputil.DelCookie(c.Writer, c.Request, "user_sess")
	httputil.DelCookie(c.Writer, c.Request, "user_last")
	c.Redirect(303, "/login")
}

func GetLoginToken(c *gin.Context) {
	in := &tokenPayload{}
	err := c.Bind(in)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := Auth(c, in.Access, in.Refresh)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// construct a login user into cache.
	user.Hash = base32.StdEncoding.EncodeToString(
		securecookie.GenerateRandomKey(32),
	)

	exp := time.Now().Add(time.Hour * 72).Unix()
	token := token.New(token.SessToken, user.Login)
	tokenstr, err := token.SignExpires(user.Hash, exp)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// update userLogin in Context.
	user.Token = tokenstr
	user.Expiry = exp
	cache := c.MustGet("cache").(cache.Cache)
	cache.SetUserLogin(user)

	c.JSON(http.StatusOK, &tokenPayload{
		Access:  tokenstr,
		Expires: exp - time.Now().Unix(),
	})
}

type tokenPayload struct {
	Access  string `json:"access_token,omitempty"`
	Refresh string `json:"refresh_token,omitempty"`
	Expires int64  `json:"expires_in,omitempty"`
}

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
