package cmd

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Dataman-Cloud/baker/client"
	"github.com/Dataman-Cloud/baker/util"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	appFile = "app.zip"
)

var BuildpackImportCmd = cli.Command{
	Name:  "import",
	Usage: "import appfiles into baker fileserver.",
	Action: func(c *cli.Context) {
		if err := buildpackImport(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "app name",
		},
		cli.StringFlag{
			Name:  "from",
			Usage: "base image",
		},
		cli.StringFlag{
			Name:  "binaryFile",
			Usage: "binary file(zip file)",
		},
		cli.StringFlag{
			Name:  "binaryPath",
			Usage: "container path of binary file.",
		},
		cli.StringFlag{
			Name:  "startupFile",
			Usage: "startup script file",
		},
		cli.StringFlag{
			Name:  "startCmd",
			Usage: "startup command",
		},
		cli.BoolFlag{
			Name:  "disconf",
			Usage: "disconf switch on-off",
		},
	},
}

func buildpackImport(c *cli.Context) error {
	// validation
	appName := c.String("name")
	if appName == "" {
		logrus.Fatal("no name in input.")
		return errors.New("no name in input.")
	}
	baseImage := c.String("from")
	if baseImage == "" {
		logrus.Fatal("no from in input.")
		return errors.New("no from in input.")
	}
	binaryFile := c.String("binaryFile")
	if binaryFile == "" {
		logrus.Fatal("no binaryFile in input.")
		return errors.New("no binaryFile in input.")
	}
	if strings.Index(binaryFile, ".zip") < 0 {
		logrus.Fatal("binaryFile is not zip file.")
		return errors.New("binaryFile is not zip file.")
	}
	binaryPath := c.String("binaryPath")
	if binaryPath == "" {
		logrus.Fatal("no binaryPath in input.")
		return errors.New("no binaryPath in input.")
	}
	startCmd := c.String("startCmd")
	startupFile := c.String("startupFile")
	disconf := c.Bool("disconf")

	// login baker server
	baseUri := c.GlobalString("server")
	client, err := client.NewHttpClient(baseUri, c.GlobalString("username"), c.GlobalString("password"))
	if err != nil {
		logrus.Fatalf("erro login baker server", err)
		return err
	}

	// upload appfiles to baker server.
	zipw, err := os.Create(appFile)
	if err != nil {
		logrus.Fatalf("error create zip file.")
		return err
	}
	defer func() {
		// remove.
		err = os.Remove(appFile)
		if err != nil {
			logrus.Error("error remove app.zip in the path.")
			return
		}
	}()

	//buf := new(bytes.Buffer)
	//w := zip.NewWriter(buf) // Create a new zip archive.
	w := zip.NewWriter(zipw)
	// Add some files to the archive.
	var files []string
	files = append(files, binaryFile)
	if startupFile != "" {
		files = append(files, startupFile)
	}
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			logrus.Fatal(err)
		}
		fileBody, err := ioutil.ReadAll(f)
		if err != nil {
			logrus.Fatal(err)
		}

		fw, err := w.Create(file)
		if err != nil {
			logrus.Fatal(err)
		}
		_, err = fw.Write([]byte(fileBody))
		if err != nil {
			logrus.Fatal(err)
		}
	}
	err = w.Close()
	if err != nil {
		logrus.Fatal(err)
	}

	extraParams := map[string]string{
		"app-name":             appName,
		"base-image":           baseImage,
		"binary-file":          binaryFile,
		"binary-path":          binaryPath,
		"start-cmd":            startCmd,
		"startup-file":         startupFile,
		"disconf-switch-onoff": strconv.FormatBool(disconf),
		"timestamp":            strconv.FormatInt(time.Now().Unix(), 10),
	}
	req, err := util.FileUploadRequest("http://"+baseUri+"/api/buildpack/import", client.Token, "uploadfile", appFile, extraParams)
	if err != nil {
		logrus.Fatalf("error create buildpack import request: %s ", err)
		return err
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Fatalf("error buildpack import request: %s", err)
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("error read response in buildpack import: %s", err)
		return err
	}
	logrus.Info(resp.Status)
	logrus.Infof(string(respBody))
	return nil
}
