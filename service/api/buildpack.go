package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Dataman-Cloud/baker/config"
	"github.com/Dataman-Cloud/baker/util"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/Dataman-Cloud/baker/executor"
)

const (
	appfilesDir = baseDir + "/appfiles"
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
	disconf, _ := strconv.ParseBool(c.Request.FormValue("disconf-switch-onoff"))

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
	desc := appfilesDir + "/" + appName + "/" + timestamp
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
	// remove.
	err = os.Remove(desc + "/" + handler.Filename)
	if err != nil {
		logrus.Error("error remove app.zip in the path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// dockerfile
	dockerfilePath := appfilesDir + "/" + appName + "/Dockerfile"
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		// create Dockerfile for the app.
		dockerfile := "FROM " + baseImage + "\n" +
			"MAINTAINER admin@dataman-inc.com" + "\n\n" +
			"#####ADD BINARY FILE########\n" +
			"ADD " + binaryFile + " " + binaryPath + "\n\n"
		if startupFile != "" {
			dockerfile += "#####COPY STARTUP FILE########\n" +
				"COPY " + startupFile + " /" + "\n\n"
		}
		// download config files.
		if disconf {
			dockerfile += "#####DOWNLOAD CONFIG FILES########\n" +
				"COPY run.sh /\n"
			if startCmd != "" {
				// copy run.sh to app timestamp directory.
				err = util.CopyFile(baseDir+"/bin/run.sh", appfilesDir+"/"+appName+"/run.sh")
				if err != nil {
					logrus.Fatal("error copy run.sh to the path.")
					c.AbortWithError(http.StatusBadRequest, err)
					return
				}
				f, err := os.OpenFile(appfilesDir+"/"+appName+"/run.sh", os.O_APPEND|os.O_WRONLY, 0777)
				if err != nil {
					logrus.Error("error open run.sh")
					c.AbortWithError(http.StatusBadRequest, err)
					return
				}
				defer f.Close()
				if _, err = f.WriteString(startCmd); err != nil {
					logrus.Error("error write startCmd string to run.sh")
					c.AbortWithError(http.StatusBadRequest, err)
					return
				}
			}
			dockerfile += "ENTRYPOINT run.sh && /bin/bash"
		}

		if !disconf {
			if startCmd != "" {
				dockerfile += "ENTRYPOINT " + startCmd + "&& /bin/bash"
			}
		}

		if startCmd != "" {
			dockerfile += "ENTRYPOINT run.sh && /bin/bash"
		}
		ioutil.WriteFile(dockerfilePath, []byte(dockerfile), 0777)
	}
	c.JSON(http.StatusOK, struct{ Filepath string }{desc})
}

// BuildpackList is a endpoint that
// list app files to baker fileserver.
func BuildpackList(c *gin.Context) {
	searchDir := appfilesDir + "" + c.Query("path")
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
	path := appfilesDir + "" + c.Query("path")
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
	appName := appfilesDir + "/" + c.Query("name") + "/Dockerfile"
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
	appName := c.Request.FormValue("app-name")
	path := appfilesDir + "/" + appName

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
	}{path})
}

// BuildpackImagePush is a endpoint that
// push docker image to docker registry.
func BuildpackImagePush(c *gin.Context) {
	appName := c.Query("name")
	timestamp := c.Query("timestamp")
	dockerfile := appfilesDir + "/" + "Dockerfile"
	path := appfilesDir + "/" + timestamp
	// copy baker to app timestamp directory.
	err := util.CopyFile(baseDir+"/bin/baker", path+"/baker")
	if err != nil {
		logrus.Fatal("error copy baker to the path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// copy run.sh to app timestamp directory.
	err = util.CopyFile(baseDir+"/bin/run.sh", path+"/run.sh")
	if err != nil {
		logrus.Fatal("error copy run.sh to the path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// copy Dockerfile to app timestamp directory.
	err = util.CopyFile(dockerfile, path+"/Dockerfile")
	if err != nil {
		logrus.Fatal("error copy dockerfile to the path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// execute an executer to imagepush.
	config := c.MustGet("config").(*config.Config)
	imageName := appName + ":" + timestamp
	workPool := config.WorkPool
	imagePushWorkPoolOptions := workPool["imagepush"]
	imagePushWorkPoolSize, _ := strconv.Atoi(imagePushWorkPoolOptions.MaxWorkers)
	works := make([]func(), 1)
	works[0] = func() {
		executor.ImagePush(imageName, &config.DockerRegistry)
	}
	executor, err := executor.NewExecutor(imagePushWorkPoolSize, works)
	executor.Execute()

	// remove.
	defer func() {
		err = os.Remove(path + "/baker")
		if err != nil {
			logrus.Error("error remove baker in the path.")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		err = os.Remove(path + "/run.sh")
		if err != nil {
			logrus.Error("error remove run.sh in the path.")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		err = os.Remove(path + "/Dockerfile")
		if err != nil {
			logrus.Error("error remove Dockerfile in the path.")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}()
}
