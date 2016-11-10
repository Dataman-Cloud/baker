package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Dataman-Cloud/baker/service/api"
	"github.com/Dataman-Cloud/baker/service/router/middleware/session"
	"github.com/Dataman-Cloud/baker/service/router/middleware/token"
)

// Load loads the router
func Load(middleware ...gin.HandlerFunc) http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	fs := http.FileServer(http.Dir("/fileserver"))
	e.GET("/static/*filepath", func(c *gin.Context) {
		fs.ServeHTTP(c.Writer, c.Request)
	})

	e.Use(middleware...)
	e.Use(session.SetUser())
	e.Use(token.Refresh)

	e.GET("/login", api.ShowLogin)
	e.GET("/logout", api.GetLogout)

	auth := e.Group("/authorize")
	{
		auth.GET("", api.GetLogin)
		auth.POST("", api.GetLogin)
		auth.POST("/token", api.GetLoginToken)
	}

	disconf := e.Group("/api/disconf")
	{
		disconf.Use(session.MustAdmin())
		disconf.POST("/push", api.DisConfPush)
		disconf.GET("/pull", api.DisConfPull)
		disconf.GET("/search", api.DisConfList)
		disconf.DELETE("/delete", api.DisConfDel)
	}

	buildpack := e.Group("/api/buildpack")
	{
		buildpack.Use(session.MustAdmin())
		buildpack.POST("/import", api.BuildpackImport)
		buildpack.GET("/search", api.BuildpackList)
		buildpack.DELETE("/delete", api.BuildpackDel)
		buildpack.POST("/dockerfile/push", api.BuildpackDockerfilePush)
		buildpack.GET("/dockerfile/pull", api.BuildpackDockerfilePull)
		buildpack.POST("/image/push", api.BuildpackImagePush)
	}
	return e
}
