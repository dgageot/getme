package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	http_headers "github.com/dgageot/getme/headers"
)

var ReleaseURL = regexp.MustCompile(`https://github.com/([^/]*)/([^/]*)/releases/download/([^/]*)/(.*)`)

type release struct {
	Assets []asset `json:"assets"`
}

type asset struct {
	Id                 int64  `json:"id"`
	BrowserDownloadURL string `json:"browser_download_url"`
	URL                string `json:"url"`
}

func AssetUrl(url string, headers []string) (string, error) {
	parts := ReleaseURL.FindStringSubmatch(url)
	org := parts[1]
	project := parts[2]
	tag := parts[3]
	assetsUrl := "https://api.github.com/repos/" + org + "/" + project + "/releases/tags/" + tag

	req, err := http.NewRequest("GET", assetsUrl, nil)
	if err != nil {
		return "", err
	}

	if err := http_headers.Add(headers, req); err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	rel := release{}
	if err := json.Unmarshal(body, &rel); err != nil {
		return "", err
	}

	for _, relAsset := range rel.Assets {
		if relAsset.BrowserDownloadURL == url {
			return relAsset.URL, nil
		}
	}

	return "", fmt.Errorf("Unable to find this release: %s", url)
}
