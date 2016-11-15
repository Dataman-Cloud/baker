package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/Dataman-Cloud/baker/client"
)

var BuildpackDelCmd = cli.Command{
	Name:  "del",
	Usage: "delete app files in fileserver.",
	Action: func(c *cli.Context) {
		if err := buildpackDel(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "delete files in the path",
		},
	},
}

func buildpackDel(c *cli.Context) error {
	path := c.String("path")
	if path == "" {
		logrus.Fatal("no path in input")
		return errors.New("no path in input")
	}

	// login baker server
	baseUri := c.GlobalString("server")
	client, err := client.NewHttpClient(baseUri, c.GlobalString("username"), c.GlobalString("password"))
	if err != nil {
		logrus.Fatalf("erro login baker server: %s", err)
		return err
	}

	// appfiles del
	logrus.Infof("delete app files in the path: %s", path)
	uri := "http://" + baseUri + "/api/buildpack/delete?path=" + path
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		logrus.Fatalf("error create buildpack delete request: %s", err)
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Token))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Fatalf("error buildpack delete request: %s", err)
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("error read response in buildpack delete: %s", err)
		return err
	}
	logrus.Info(resp.Status)
	logrus.Infof(string(respBody))
	return nil
	return nil
}
