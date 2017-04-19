package files

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgageot/getme/github"
	"github.com/pkg/errors"
)

// Download downloads an url to a destination file. Additional headers can be given.
// This is helpful to pass authentication tokens.
func Download(url string, destination string, headers []string) error {
	destinationTmp := destination + ".tmp"

	err := downloadURL(url, destinationTmp, headers)
	if err == nil {
		return nil
	}

	if !github.ReleaseURL.MatchString(url) {
		return err
	}

	log.Println("Github release url detected")

	assetUrl, err := github.AssetUrl(url, headers)
	if err != nil {
		return err
	}

	log.Println("Github asset url is:", assetUrl)

	headers = append(headers, "Accept=application/octet-stream")
	err = downloadURL(assetUrl, destinationTmp, headers)
	if err != nil {
		return err
	}

	return os.Rename(destinationTmp, destination)
}

func downloadURL(url string, destination string, headers []string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if err := addHeaders(headers, req); err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return errors.New(resp.Status)
	}

	return CopyFrom(destination, 0666, resp.Body)
}

func addHeaders(headers []string, req *http.Request) error {
	for _, header := range headers {
		parts := strings.Split(header, "=")
		if len(parts) != 2 {
			return fmt.Errorf("Invalid header [%s]. Should be [key=value]", header)
		}
		req.Header.Add(parts[0], parts[1])
	}

	return nil
}
