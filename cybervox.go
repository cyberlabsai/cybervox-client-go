// Package cybervox implements all structs and functions necessary to talk to the CyberVox api platform.
package cybervox

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
)

var log = logrus.WithField("package", "cybervox")

func getAccessToken(clientID, clientSecret string) (accessToken string, err error) {
	if accessToken = getCachedAccessToken(); accessToken != "" {
		log.Debugln("using cached access token...")
		return
	}

	var response oauthResponse
	if response, err = fetchAccessToken(clientID, clientSecret); err != nil {
		return
	}
	// go ahead even if we cannot save it
	_ = saveAccessToken(response.AccessToken)

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
