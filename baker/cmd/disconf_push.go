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

	"github.com/Dataman-Cloud/baker/client"
)

const (
	serverUrl = "127.0.0.1:8000"
	username  = "admin"
	password  = "badmin"
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
		cli.StringFlag{
			Name:  "serverUrl",
			Usage: "baker server uri",
			Value: serverUrl,
		},
		cli.StringFlag{
			Name:  "username",
			Usage: "username to login baker server.",
			Value: username,
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "password to login baker server.",
			Value: password,
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
	logrus.Infof("push config files into disconf.")
	var err error
	ymlfile := c.String("ymlfile")
	buf, err := ioutil.ReadFile(ymlfile)
	if err != nil {
		logrus.Fatalf("error open ymlfile: %s", ymlfile)
		return err
	}
	upload := &Upload{}
	err = yaml.Unmarshal(buf, upload)
	if err != nil {
		logrus.Fatalf("error unmarshal jsonfile:%s", ymlfile)
		return err
	}

	baseUrl := c.String("serverUrl")
	client, err := client.NewHttpClient(baseUrl, c.String("username"), c.String("password"))
	if err != nil {
		logrus.Fatalf("erro login baker server", err)
		return err
	}

	logrus.Infof("upload files size:%d", len(upload.Files))
	for filename, fileparam := range upload.Files {
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
			logrus.Fatalf("error open file: %s", filename)
			return err
		}

		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			logrus.Fatalf("error copy file: %s", filename)
			return err
		}

		// post upload request to baker server.
		uploadUrl := fmt.Sprintf("http://%s/api/disconf/push?download-path=%s&container-path=%s",
			baseUrl,
			url.QueryEscape(fileparam.DownloadPath),
			url.QueryEscape(fileparam.ContainerPath))
		contentType := bodyWriter.FormDataContentType()

		req, err := http.NewRequest("POST", uploadUrl, bodyBuf)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Token))
		req.Header.Set("Content-Type", contentType)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Fatal("error upload to baker server.")
			return err
		}
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Fatal("error read response in upload request.")
			return err
		}
		logrus.Info(resp.Status)
		logrus.Infof(string(respBody))
	}
	return nil
}
