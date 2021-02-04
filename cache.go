package cybervox

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
)

var digitsRegexp = regexp.MustCompile(`\d+`)

func getCachedAccessToken() (accessToken string) {
	var err error
	var matches []string
	if matches, err = filepath.Glob("access-token-*.txt"); err != nil || len(matches) == 0 {
		return
	}
	log.Debugf("found %d access token files", len(matches))

	//
	// extract the file creation time from its name
	//
	cacheFilename := matches[len(matches)-1]
	creationTime := digitsRegexp.FindString(cacheFilename)

	var unix int64
	if unix, err = strconv.ParseInt(creationTime, 10, 64); err != nil {
		return
	}

	//
	// check for expired cache
	//
	cachedFor := time.Since(time.Unix(unix, 0)).Hours()
	log.Debugf("cache created %.2f hours ago", cachedFor)
	if cachedFor > 23.5 {
		// cleanup and leave
		for _, filename := range matches {
			_ = os.Remove(filename)
		}
		return
	}

	//
	// still valid, use it
	//
	log.Debugf("reading access token from %q", cacheFilename)
	var contents []byte
	if contents, err = ioutil.ReadFile(cacheFilename); err != nil {
		return
	}

	return string(contents)
}

func saveAccessToken(token string) error {
	cacheFilename := fmt.Sprintf("access-token-%d.txt", time.Now().Unix())
	log.Debugf("saving access token to %q", cacheFilename)
	return ioutil.WriteFile(cacheFilename, []byte(token), 0644)
}

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

	log.Debugln("fetching new access token...")
	res, _, errs := gorequest.New().Post("https://cybervox-dev.us.auth0.com/oauth/token").
		Send(request).
		EndStruct(&response)
	if errs != nil {
		return oauthResponse{}, fmt.Errorf("gorequest: %v", errs)
	}
	if res.StatusCode != http.StatusOK {
		return oauthResponse{}, fmt.Errorf("http.Status: %v", res.Status)
	}
	return response, nil
}
