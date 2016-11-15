package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Dataman-Cloud/baker/util"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// BuildpackImport is a endpoint that
// import app files to baker fileserver.
func BuildpackImport(c *gin.Context) {
	// get request params.
	appName := c.Request.FormValue("app-name")
	baseImage := c.Request.FormValue("base-image")
	binaryFile := c.Request.FormValue("binary-file")
	binaryPath := c.Request.FormValue("binary-path")
	startCmd := c.Request.FormValue("start-cmd")
	startupFile := c.Request.FormValue("startup-file")
	timestamp := c.Request.FormValue("timestamp")

	// parse upload file.
	c.Request.ParseMultipartForm(32 << 20)
	file, handler, err := c.Request.FormFile("uploadfile")
	if err != nil {
		logrus.Error("error parse upload file.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	// save upload file to desc dir.
	desc := baseDir + "/appfiles/" + appName + "/" + timestamp
	err = os.MkdirAll(desc, 0777)
	if err != nil {
		logrus.Error("error create desc file directory.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	f, err := os.OpenFile(desc+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Error("error save zip file. ")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	// unzip upload file in desc dir.
	err = util.Unzip(desc+"/"+handler.Filename, desc)
	if err != nil {
		logrus.Error("error unzip app file.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// dockerfile
	dockerfilePath := baseDir + "/appfiles/" + appName + "/Dockerfile"
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		// create Dockerfile for the app.
		dockerfile := "FROM " + baseImage + "\n\n" +
			"MAINTAIN admin@dataman-inc.com" + "\n\n" +
			"ADD " + binaryFile + " " + binaryPath + "\n"
		if startupFile != "" {
			dockerfile += "COPY " + startupFile + " /" + "\n\n"
		}
		// add bakercli pull properties file from disconf depending on envVars.
		bakercli := baseDir + "/bin/baker"
		dockerfile += "# DOWNLOAD PROPERTY FILES FROM DISCONF IN BAKER SERVER\n" +
			"COPY  " + bakercli + " /\n" +
			"RUN ./baker disconf pull --path=$CONFIG_DIR && \n" +
			"./baker disconf unzip --file=props.zip --path=/" + appName + " &&\n" +
			"mv $CONFIG_DIR /\n"
		if startCmd != "" {
			dockerfile += "CMD [\"" + startCmd + "\"]"
		}
		ioutil.WriteFile(dockerfilePath, []byte(dockerfile), 0777)
	}
	c.JSON(http.StatusOK, struct{ Filepath string }{desc})
}

// BuildpackList is a endpoint that
// list app files to baker fileserver.
func BuildpackList(c *gin.Context) {
	searchDir := baseDir + "/appfiles" + c.Query("path")
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

// BuildpackDel is a endpoint that
// delete app files to baker fileserver.
func BuildpackDel(c *gin.Context) {
	path := baseDir + "/appfiles" + c.Query("path")
	err := os.RemoveAll(path)
	if err != nil {
		logrus.Error("error remove files in the path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, struct {
		Message string
	}{"appfiles del is ok."})
}

// BuildpackDockerfilePull is a endpoint that
// pull dockerfile from baker fileserver.
func BuildpackDockerfilePull(c *gin.Context) {
	writer := c.Writer
	appName := baseDir + "/appfiles/" + c.Query("name") + "/Dockerfile"
	// write file content to response body
	openFile, err := os.Open(appName)
	if err != nil {
		logrus.Error("error open app dockerfile.")
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
	writer.Header().Set("Content-Disposition", "attachment; filename="+appName)
	writer.Header().Set("Content-Type", fileContentType)
	writer.Header().Set("Content-Length", fileSize)

	//Send the file
	//We read 512 bytes from the file already so we reset the offset back to 0
	openFile.Seek(0, 0)
	io.Copy(writer, openFile) //'Copy' the file to the client
}

// BuildpackDockerfilePush is a endpoint that
// push dockerfile to baker fileserver.
func BuildpackDockerfilePush(c *gin.Context) {

}

// BuildpackImagePush is a endpoint that
// push docker image to docker registry.
func BuildpackImagePush(c *gin.Context) {

}
