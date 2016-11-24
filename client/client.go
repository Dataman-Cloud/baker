package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type Client struct {
	Username string
	Token    string
}

func NewHttpClient(uri, username, password string) (*Client, error) {
	// login request to baker server.
	loginUrl := "http://" + uri + "/authorize"
	loginPayLoad := "{\"username\":\"" + username + "\",\"password\":\"" + password + "\"}"

	req, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer([]byte(loginPayLoad)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("erro login.")
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("error read response in login request.")
		return nil, err
	}
	logrus.Info(resp.Status)
	logrus.Infof(string(respBody))
	type Login struct {
		Access  string `json:"access_token"`
		Refresh int64  `json:"expires_in"`
	}
	var login Login
	err = json.Unmarshal(respBody, &login)
	if err != nil {
		logrus.Error("error unmarshal for login response.%s", err)
		return nil, err
	}

	return &Client{
		Username: username,
		Token:    login.Access,
	}, nil
}
