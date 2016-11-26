package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/Dataman-Cloud/baker/config"
	"github.com/Dataman-Cloud/baker/util"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/Dataman-Cloud/baker/executor"
)

const (
	appfilesDir  = baseDir + "/appfiles"
	workspaceDir = baseDir + "/workspace"
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

	// create appDir.
	appDir := appfilesDir + "/" + appName
	desc := appDir + "/" + timestamp
	err = os.MkdirAll(desc, 0777)
	if err != nil {
		logrus.Error("error create desc file directory.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// save app zip file to desc dir.
	f, err := os.OpenFile(desc+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Error("error save zip file. ")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	// unzip app zip file in desc dir.
	err = util.Unzip(desc+"/"+handler.Filename, desc)
	if err != nil {
		logrus.Error("error unzip app file.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// remove.
	defer func() {
		err = os.Remove(desc + "/" + handler.Filename)
		if err != nil {
			logrus.Error("error remove app.zip in the path.")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}()
	// dockerfile
	dockerfilePath := appDir + "/" + "Dockerfile"
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		// create Dockerfile for the app.
		dockerfile := "FROM " + baseImage + "\n" +
			"MAINTAINER admin@dataman-inc.com" + "\n\n" +
			"#####ADD BINARY FILE########\n" +
			"ADD " + binaryFile + " " + binaryPath + "\n\n" +
			"WORKDIR /\n"
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
					logrus.Error("error copy run.sh to the path.")
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
			dockerfile += "ENTRYPOINT run.sh && /bin/bash \n"
		}

		if !disconf {
			if startCmd != "" {
				dockerfile += "ENTRYPOINT " + startCmd + "&& /bin/bash"
			}
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
	appName := c.Query("name")
	appDir := appfilesDir + "/" + appName

	// write file content to response body
	dockerfile := appDir + "/" + "Dockerfile"
	openFile, err := os.Open(dockerfile)
	if err != nil {
		logrus.Error("error open app dockerfile.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer openFile.Close()
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
	writer.Header().Set("Content-Disposition", "attachment; filename="+dockerfile)
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
	appDir := appfilesDir + "/" + appName

	c.Request.ParseMultipartForm(32 << 20)
	file, handler, err := c.Request.FormFile("uploadfile")
	if err != nil {
		logrus.Error("error get upload file.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	// appDir lock.
	d, err := os.Open(appDir)
	if err != nil {
		logrus.Errorf("error open app path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = syscall.Flock(int(d.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		logrus.Errorf("cannot flock app path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// unlock.
	defer func() {
		d.Close()
		syscall.Flock(int(d.Fd()), syscall.LOCK_UN)
	}()

	// write file.
	f, err := os.OpenFile(appDir+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Error("error create upload file. ")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	c.JSON(http.StatusOK, struct {
		Filepath string
	}{appDir})
}

// BuildpackImagePush is a endpoint that
// push docker image to docker registry.
func BuildpackImagePush(c *gin.Context) {
	appName := c.Query("name")
	timestamp := c.Query("timestamp")

	appDir := appfilesDir + "/" + appName
	taskID := strconv.FormatInt(time.Now().Unix(), 10)
	workDir := workspaceDir + "/" + taskID

	// setup imagepush workspace directory.
	err := setupImagePushWorkDir(appDir, timestamp, workDir)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	//create an executer to imagepush.
	cf := c.MustGet("config").(*config.Config)
	bakeWorkPool := c.MustGet("bakeworkpool").(*executor.WorkPool)
	imageName := appName + ":" + timestamp
	taskStats := make(chan int)
	taskMsg := make(chan string)
	isDone := make(chan bool)
	imagePushTask := executor.NewImagePushTask(workDir, imageName, &cf.DockerRegistry)
	taskCollector := executor.NewCollector(taskID, taskStats, taskMsg, isDone)
	work := createImagePushWork(imagePushTask, taskCollector)

	tasks := make([]*executor.Task, 1)
	tasks[0] = &executor.Task{
		ID:   taskID,
		Work: work,
	}
	taskExec, err := executor.NewExecutor(bakeWorkPool, tasks, taskCollector)
	if err != nil {
		logrus.Error("error create job executor.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	taskExec.Execute()

	c.JSON(http.StatusOK, struct {
		WorkDir string
	}{workDir})

}

// copy app files, baker, run.sh, dockerfile to workspace directory
func setupImagePushWorkDir(appDir, timestamp, workDir string) error {
	runshFile := appDir + "/" + "run.sh"
	dockerFile := appDir + "/" + "Dockerfile"
	path := appDir + "/" + timestamp
	// copy appfiles to workspace directory.
	var err error
	err = util.CopyDir(path, workDir)
	if err != nil {
		logrus.Error("error copy app files to the path.")
		return err
	}
	// copy baker to workspace directory.
	err = util.CopyFile(baseDir+"/bin/baker", workDir+"/baker")
	if err != nil {
		logrus.Error("error copy baker to the path.")
		return err
	}
	// copy run.sh to workspace directory.
	err = util.CopyFile(runshFile, workDir+"/run.sh")
	if err != nil {
		logrus.Error("error copy run.sh to the path.")
		return err
	}
	// copy Dockerfile to workspace directory.
	err = util.CopyFile(dockerFile, workDir+"/Dockerfile")
	if err != nil {
		logrus.Error("error copy dockerfile to the path.")
		return err
	}
	return nil
}

// create ImagePushWork
func createImagePushWork(t *executor.ImagePushTask, c *executor.Collector) func() {
	return func() {
		c.TaskStats <- executor.StatusRunning
		var err error
		workDir := t.WorkDir
		imageName := t.ImageName
		dockerRegistry := t.Config
		imagePushTask := executor.NewImagePushTask(workDir, imageName, dockerRegistry)

		// dockerLogin
		err = imagePushTask.DockerLogin()
		if err != nil {
			logrus.Error("error execute docker login.")
			c.TaskMsg <- err.Error()
			return
		}

		// dockerBuild
		err = imagePushTask.DockerBuild()
		if err != nil {
			logrus.Error("error execute docker build.")
			c.TaskMsg <- err.Error()
			return
		}

		// dockerPush
		err = imagePushTask.DockerPush()
		if err != nil {
			logrus.Error("error execute docker push.")
			c.TaskMsg <- err.Error()
			return
		}

		defer func() {
			if r := recover(); r != nil {
				c.TaskStats <- executor.StatusFailed
			} else {
				c.TaskStats <- executor.StatusFinished
			}
		}()
	}
}
