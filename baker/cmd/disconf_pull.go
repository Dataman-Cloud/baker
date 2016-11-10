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

var DisConfPullCmd = cli.Command{
	Name:  "pull",
	Usage: "pull config files from disconfig",
	Action: func(c *cli.Context) {
		if err := disConfPull(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "pull config files in path",
		},
	},
}

func disConfPull(c *cli.Context) error {
	path := c.String("path")
	if path == "" {
		logrus.Fatal("no path in input")
		return errors.New("no path in input")
	}

	// login baker server
	baseUri := c.GlobalString("server")
	client, err := client.NewHttpClient(baseUri, c.GlobalString("username"), c.GlobalString("password"))
	if err != nil {
		logrus.Fatalf("erro login baker server", err)
		return err
	}

	// disconf pull
	logrus.Infof("download config files in the path: %s", path)
	uri := "http://" + baseUri + "/api/disconf/pull?path=" + path
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		logrus.Fatalf("error new disconf pull request: %s", err)
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Token))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Fatalf("error disconf pull request: %s", err)
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

	zipfile := "props.zip"
	out, err := os.Create(zipfile)
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
