package api

import (
	"io"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type disConfPayload struct {
	Filepath string
}

// DisConfPush is a endpoint that
// push config files defined in jsonfile into config management.
func DisConfPush(c *gin.Context) {
	appName := c.Request.FormValue("app-name")
	label := c.Request.FormValue("label")
	containerPath := c.Request.FormValue("container-path")
	c.Request.ParseMultipartForm(32 << 20)
	file, handler, err := c.Request.FormFile("uploadfile")
	if err != nil {
		logrus.Error("error get upload file.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer file.Close()
	err = os.MkdirAll("/fileserver/disconf"+"/"+appName+"/"+label+""+containerPath, 0777)
	if err != nil {
		logrus.Error("error create upload file directory")
		return
	}
	f, err := os.OpenFile("/fileserver/disconf"+"/"+appName+"/"+label+""+containerPath+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Error("error upload file. ")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	c.JSON(http.StatusOK, &disConfPayload{Filepath: containerPath})
}

// DisConfPull is a endpoint that
// pull config files in config management.
func DisConfPull(c *gin.Context) {

}

// DisConfList is a endpoint that
// list config files in config management.
func DisConfList(c *gin.Context) {

}

// DisConfDel is a endpoint that
// delete config files in config management.
func DisConfDel(c *gin.Context) {

}
