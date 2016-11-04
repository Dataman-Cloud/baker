package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	serverUrl = "127.0.0.1:8000"
)

var DisConfPushCmd = cli.Command{
	Name:  "push",
	Usage: "push config files into config management",
	Action: func(c *cli.Context) {
		if err := disConfPush(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "ymlfile",
			Usage: "ymlfile for uploading config files",
			Value: "push.yml",
		},
	},
}

type Upload struct {
	Files map[string]FileParam `yaml:"uploads"`
}

type FileParam struct {
	DownloadPath  string `yaml:"download_path"`  // download path.
	ContainerPath string `yaml:"container_path"` // container path.
}

func disConfPush(c *cli.Context) error {
	logrus.Infof("push config files into management.")
	ymlfile := c.String("ymlfile")
	buf, err := ioutil.ReadFile(ymlfile)
	if err != nil {
		logrus.Fatalf("Fail to open ymlfile: %s", ymlfile)
		return err
	}
	upload := &Upload{}
	err = yaml.Unmarshal(buf, upload)
	if err != nil {
		logrus.Fatalf("Fail to Unmarshal jsonfile:%s", ymlfile)
		return err
	}
	logrus.Infof("upload files size:%d", len(upload.Files))
	for filename, fileparam := range upload.Files {
		logrus.Infof("config filename:%s", filename)

		// simulate web client to upload file.
		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)
		defer bodyWriter.Close()
		fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
		if err != nil {
			logrus.Fatalf("error write to buffer: %s", filename)
			return err
		}

		fh, err := os.Open(filename)
		if err != nil {
			logrus.Fatalf("error open file : %s", filename)
			return err
		}

		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			logrus.Fatalf("error copy file: %s", filename)
			return err
		}
		contentType := bodyWriter.FormDataContentType()
		logrus.Infof("content-type:%s", contentType)

		// post upload request to baker server.
		downloadPath := url.QueryEscape(fileparam.DownloadPath)
		containerPath := url.QueryEscape(fileparam.ContainerPath)
		uploadUrl := fmt.Sprintf("http://%s/api/disconf/push?download-path=%s&container-path=%s", serverUrl, downloadPath, containerPath)
		logrus.Infof("uploadUrl:%s", uploadUrl)
		resp, err := http.Post(uploadUrl, contentType, bodyBuf)
		if err != nil {
			logrus.Fatal("error post upload request to baker server.")
			return err
		}
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Fatal("error read response body.")
			return err
		}
		logrus.Info(resp.Status)
		logrus.Infof(string(resp_body))
	}
	return nil
}
