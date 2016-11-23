package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/Dataman-Cloud/baker/util"
)

const (
	disconfDir = baseDir + "/disconf"
)

// DisConfPush is a endpoint that
// push config files in disconf.
func DisConfPush(c *gin.Context) {
	appName := c.Request.FormValue("app-name")
	label := c.Request.FormValue("label")
	timestamp := c.Request.FormValue("timestamp")
	containerPath := c.Request.FormValue("container-path")
	path := disconfDir + "/" + appName + "/" + label + "/" + timestamp + "" + containerPath

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
	c.JSON(http.StatusOK, struct {
		Filepath string
	}{containerPath})
}

// DisConfPull is a endpoint that
// pull config files in disconfig.
func DisConfPull(c *gin.Context) {
	writer := c.Writer
	path := disconfDir + "" + c.Query("path")
	zipfile := tmpDir + "/" + "props.zip"
	err := util.Zipit(path, zipfile)
	if err != nil {
		logrus.Error("error zip prop file. ")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// remove.
	defer func() {
		err = os.Remove(zipfile)
		if err != nil {
			logrus.Error("error remove props.zip in tmp path.")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}()

	// write zipfile content to response body
	openFile, err := os.Open(zipfile)
	if err != nil {
		logrus.Error("error open zipfile.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	openFile.Read(fileHeader)
	//Get content type of file
	fileContentType := http.DetectContentType(fileHeader)

	//Get the file size
	fileStat, _ := openFile.Stat()                     //Get info from file
	fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string

	//Send the headers
	writer.Header().Set("Content-Disposition", "attachment; filename="+zipfile)
	writer.Header().Set("Content-Type", fileContentType)
	writer.Header().Set("Content-Length", fileSize)

	//Send the file
	//We read 512 bytes from the file already so we reset the offset back to 0
	openFile.Seek(0, 0)
	io.Copy(writer, openFile) //'Copy' the file to the client
}

// DisConfList is a endpoint that
// list config files in disconfig.
func DisConfList(c *gin.Context) {
	searchDir := disconfDir + "" + c.Query("path")
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
	c.JSON(http.StatusOK, struct {
		FileList []string
	}{fileList})
}

// DisConfDel is a endpoint that
// delete config files in disconfig.
func DisConfDel(c *gin.Context) {
	path := disconfDir + "" + c.Query("path")
	err := os.RemoveAll(path)
	if err != nil {
		logrus.Error("error remove all files in the path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, struct {
		Message string
	}{"disconf del is ok."})
}
