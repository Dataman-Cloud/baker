package cmd

import (
	"bytes"
	"encoding/json"
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
	username  = "admin"
	password  = "badmin"
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
	logrus.Infof("push config files into disconf.")
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
	// login request to baker server.
	loginUrl := fmt.Sprintf("http://%s/authorize", serverUrl)
	loginPayLoad := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\"}", username, password)
	loginReq, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer([]byte(loginPayLoad)))
	loginReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	loginResp, err := client.Do(loginReq)
	if err != nil {
		logrus.Fatal("erro login.")
		return err
	}
	defer loginResp.Body.Close()
	loginRespBody, err := ioutil.ReadAll(loginResp.Body)
	if err != nil {
		logrus.Fatal("error read response in login request.")
		return err
	}
	logrus.Info(loginResp.Status)
	logrus.Infof(string(loginRespBody))
	type Login struct {
		Access  string `json:"access_token"`
		Refresh int64  `json:"expires_in"`
	}
	var login Login
	err = json.Unmarshal(loginRespBody, &login)
	if err != nil {
		logrus.Fatalf("error unmarshal for login response.%s", err)
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
		downloadPath := url.QueryEscape(fileparam.DownloadPath)
		containerPath := url.QueryEscape(fileparam.ContainerPath)
		uploadUrl := fmt.Sprintf("http://%s/api/disconf/push?download-path=%s&container-path=%s", serverUrl, downloadPath, containerPath)
		contentType := bodyWriter.FormDataContentType()

		req, err := http.NewRequest("POST", uploadUrl, bodyBuf)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", login.Access))
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
