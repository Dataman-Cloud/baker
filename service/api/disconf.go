package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	baseDir = "/fileserver/disconf"
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
	path := baseDir + "/" + appName + "/" + label + "" + containerPath

	c.Request.ParseMultipartForm(32 << 20)
	file, handler, err := c.Request.FormFile("uploadfile")
	if err != nil {
		logrus.Error("error get upload file.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer file.Close()
	err = os.MkdirAll(path, 0777)
	if err != nil {
		logrus.Error("error create upload file directory.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	f, err := os.OpenFile(path+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Error("error create upload file. ")
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
	path := c.Query("path")
	logrus.Infof("path:%s", path)
}

// DisConfList is a endpoint that
// list config files in config management.
func DisConfList(c *gin.Context) {
	searchDir := baseDir + "" + c.Query("path")
	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		logrus.Error("error list file.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, fileList)
}

// DisConfDel is a endpoint that
// delete config files in config management.
func DisConfDel(c *gin.Context) {

}
