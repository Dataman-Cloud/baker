package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
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
			Name:  "name",
			Usage: "application name",
		},
		cli.StringFlag{
			Name:  "label",
			Usage: "application label name (such as: dev,test,prod....).",
		},
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
	ContainerPath string `yaml:"container_path"` // container path.
}

func disConfPush(c *cli.Context) error {
	// validation.
	appName := c.String("name")
	if appName == "" {
		logrus.Fatal("no name input")
		return errors.New("no name input.")
	}
	label := c.String("label")
	if label == "" {
		logrus.Fatal("no label input")
		return errors.New("no label input.")
	}

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
		path, _ := os.Getwd()
		path += "/" + filename
		extraParams := map[string]string{
			"app-name":       appName,
			"label":          label,
			"container-path": fileparam.ContainerPath,
		}
		req, err := fileUploadRequest("http://"+baseUrl+"/api/disconf/push", client.Token, "uploadfile", path, extraParams)
		if err != nil {
			logrus.Fatalf("error make file upload request.%s ", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Fatalf("error upload to baker server.%s", err)
			return err
		}
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Fatalf("error read response in upload request.%s", err)
			return err
		}
		logrus.Info(resp.Status)
		logrus.Infof(string(respBody))
	}
	return nil
}

// fileUploadRequest is create a file upload http request with optional extra params
func fileUploadRequest(uri string, token, uploadParamName, uploadFilePath string, extraParams map[string]string) (*http.Request, error) {
	file, err := os.Open(uploadFilePath)
	if err != nil {
		return nil, err
	}
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(uploadParamName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContent)

	for key, val := range extraParams {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	contentType := writer.FormDataContentType()

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", contentType)
	return req, err
}
