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

var BuildpackListCmd = cli.Command{
	Name:  "list",
	Usage: "list app files in fileserver.",
	Action: func(c *cli.Context) {
		if err := buildpackList(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "list app files in the path.",
		},
	},
}

func buildpackList(c *cli.Context) error {
	path := c.String("path")
	if path == "" {
		logrus.Fatal("no path input")
		return errors.New("no path input")
	}

	// login baker server
	baseUri := c.GlobalString("server")
	client, err := client.NewHttpClient(baseUri, c.GlobalString("username"), c.GlobalString("password"))
	if err != nil {
		logrus.Fatalf("erro login baker server: %s", err)
		return err
	}

	// appfiles list.
	logrus.Infof("list app files in baker server: %s", path)
	uri := "http://" + baseUri + "/api/buildpack/search?path=" + path
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		logrus.Fatalf("error create buildpack list request: %s", err)
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Token))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Fatalf("error buildpack list request: %s", err)
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("error read response in buildpack list: %s", err)
		return err
	}
	logrus.Info(resp.Status)
	logrus.Infof(string(respBody))

	return nil
}
