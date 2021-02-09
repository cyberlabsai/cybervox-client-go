// Package cybervox implements all structs and functions necessary to talk to the CyberVox api platform.
package cybervox

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

type (
	oauthRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Audience     string `json:"audience"`
		GrantType    string `json:"grant_type"`
	}
	oauthResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
)

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
)

var log = logrus.WithField("package", "cybervox")

func fetchAccessToken(clientID string, clientSecret string) (oauthResponse, error) {
	var (
		request = oauthRequest{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Audience:     "https://api.cybervox.ai",
			GrantType:    "client_credentials",
		}
		response oauthResponse
	)

	log.Debugln("fetching access token...")
	res, body, errs := gorequest.New().Post("https://api.cybervox.ai/auth").
		Send(request).
		EndStruct(&response)
	if errs != nil {
		return oauthResponse{}, fmt.Errorf("gorequest(%s): %v", string(body), errs)
	}
	if res.StatusCode != http.StatusOK {
		return oauthResponse{}, fmt.Errorf("http.Status: %v", res.Status)
	}
	return response, nil
}

func getAccessToken(clientID, clientSecret string) (string, error) {
	var response oauthResponse
	var err error
	if response, err = fetchAccessToken(clientID, clientSecret); err != nil {
		return "", err
	}

	return response.AccessToken, nil
}

// Dial connects to the API's websocket. It assumes the env vars `CLIENT_ID` and `CLIENT_SECRET` are defined correctly.
// It returns the websocket.Conn, the http.Response and maybe an error.
func Dial() (*websocket.Conn, *http.Response, error) {
	if clientID == "" || clientSecret == "" {
		return nil, nil, fmt.Errorf(`abort: check "CLIENT_ID" and "CLIENT_SECRET" envvars`)
	}
	var token string
	var err error
	if token, err = getAccessToken(clientID, clientSecret); err != nil {
		return nil, nil, err
	}
	return websocket.DefaultDialer.Dial("wss://api.cybervox.ai/ws?access_token="+token, nil)
}
