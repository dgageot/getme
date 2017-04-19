package appveyor

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"

	http_headers "github.com/dgageot/getme/headers"
)

var ArtifactURL = regexp.MustCompile(`https://ci.appveyor.com/project/([^/]*)/([^/]*)/build/([^/]*)/artifacts/(.*)`)

type info struct {
	Build build `json:"build"`
}

type build struct {
	Jobs []job `json:"jobs"`
}

type job struct {
	Id string `json:"jobId"`
}

func ArtifactUrl(url string, headers []string) (string, error) {
	parts := ArtifactURL.FindStringSubmatch(url)
	account := parts[1]
	project := parts[2]
	buildNumber := parts[3]
	artifact := parts[4]
	buildUrl := "https://ci.appveyor.com/api/projects/" + account + "/" + project + "/build/" + buildNumber

	req, err := http.NewRequest("GET", buildUrl, nil)
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

	i := info{}
	if err := json.Unmarshal(body, &i); err != nil {
		return "", err
	}

	jobId := i.Build.Jobs[0].Id

	artifactUrl := "https://ci.appveyor.com/api/buildjobs/" + jobId + "/artifacts/" + artifact

	return artifactUrl, nil
}
