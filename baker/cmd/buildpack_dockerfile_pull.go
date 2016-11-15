package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	_ "path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/Dataman-Cloud/baker/client"
)

var BuildpackDockerfilePullCmd = cli.Command{
	Name:  "pull",
	Usage: "execute pull dockerfile in fileserver",
	Action: func(c *cli.Context) {
		if err := buildpackDockerfilePull(c); err != nil {
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

func buildpackDockerfilePull(c *cli.Context) error {
	appName := c.String("name")
	if appName == "" {
		logrus.Fatal("no name input")
		return errors.New("no name input")
	}

	// login baker server
	baseUri := c.GlobalString("server")
	client, err := client.NewHttpClient(baseUri, c.GlobalString("username"), c.GlobalString("password"))
	if err != nil {
		logrus.Fatalf("erro login baker server", err)
		return err
	}

	// dockerfile pull
	logrus.Infof("download dockerfile for app: %s", appName)
	uri := "http://" + baseUri + "/api/buildpack/dockerfile/pull?name=" + appName
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		logrus.Fatalf("error create dockerfile pull request: %s", err)
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Token))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Fatalf("error dockerfile pull request: %s", err)
		return err
	}
	defer resp.Body.Close()
	//respBody, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	logrus.Fatalf("error read response in disconf pull: %s", err)
	//	return err
	//}
	//logrus.Info(resp.Status)
	//logrus.Infof(string(respBody))

	fw := "Dockerfile"
	out, err := os.Create(fw)
	if err != nil {
		logrus.Fatalf("error create download file: %s", err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logrus.Fatalf("error save file: %s", err)
		return err
	}
	return nil
}
