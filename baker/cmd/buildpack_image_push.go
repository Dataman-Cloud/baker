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

var BuildpackImagePushCmd = cli.Command{
	Name:  "push",
	Usage: "execute push image to docker registry",
	Action: func(c *cli.Context) {
		if err := buildpackImagePush(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "app name",
		},
		cli.StringFlag{
			Name:  "timestamp",
			Usage: "timestamp",
		},
	},
}

func buildpackImagePush(c *cli.Context) error {
	appName := c.String("name")
	if appName == "" {
		logrus.Fatal("no name input")
		return errors.New("no name input")
	}
	timestamp := c.String("timestamp")
	if timestamp == "" {
		logrus.Fatal("no timestamp input")
		return errors.New("no timestamp input")
	}

	// login baker server
	baseUri := c.GlobalString("server")
	client, err := client.NewHttpClient(baseUri, c.GlobalString("username"), c.GlobalString("password"))
	if err != nil {
		logrus.Fatalf("erro login baker server: %s", err)
		return err
	}
	// image push.
	logrus.Infof("app image push in baker server: %s", appName)
	uri := "http://" + baseUri + "/api/buildpack/image/push?name=" + appName + "&timestamp=" + timestamp
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		logrus.Fatalf("error create image push request: %s", err)
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Token))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Fatalf("error image push request: %s", err)
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("error read response in image push: %s", err)
		return err
	}
	logrus.Info(resp.Status)
	logrus.Infof(string(respBody))

	return nil
}
