package cmd

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/Dataman-Cloud/baker/client"
	"github.com/Dataman-Cloud/baker/util"
)

var BuildpackDockerfilePushCmd = cli.Command{
	Name:  "push",
	Usage: "execute push dockerfile in fileserver",
	Action: func(c *cli.Context) {
		if err := buildpackDockerfilePush(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "app name",
		},
	},
}

func buildpackDockerfilePush(c *cli.Context) error {
	// validation.
	appName := c.String("name")
	if appName == "" {
		logrus.Fatal("no name in input.")
		return errors.New("no name in input.")
	}

	// login baker server
	baseUri := c.GlobalString("server")
	client, err := client.NewHttpClient(baseUri, c.GlobalString("username"), c.GlobalString("password"))
	if err != nil {
		logrus.Fatalf("erro login baker server", err)
		return err
	}

	// dockerfile push
	logrus.Infof("push dockerfile to baker server.")
	path, _ := os.Getwd()
	dockerfile := path + "/" + "Dockerfile"
	extraParams := map[string]string{
		"app-name": appName,
	}
	req, err := util.FileUploadRequest("http://"+baseUri+"/api/buildpack/dockerfile/push", client.Token, "uploadfile", dockerfile, extraParams)
	if err != nil {
		logrus.Fatalf("error create dockerfile push request: %s ", err)
		return err
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Fatalf("error dockerfile push request: %s", err)
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("error read response in dockerfile push: %s", err)
		return err
	}
	logrus.Info(resp.Status)
	logrus.Infof(string(respBody))
	return nil
}
