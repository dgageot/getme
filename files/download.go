package files

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dgageot/getme/appveyor"
	"github.com/dgageot/getme/github"
	http_headers "github.com/dgageot/getme/headers"
	minio "github.com/minio/minio-go"
	"github.com/pkg/errors"
)

type Options struct {
	AuthToken            string
	AuthTokenEnvVariable string
	S3AccessKey          string
	S3SecretKey          string
	Sha256               string
}

// Download downloads an url to a destination file. Additional headers can be given.
// This is helpful to pass authentication tokens.
func Download(rawURL string, destination string, options Options) error {
	parsedUrl, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	destinationTmp := destination + ".tmp"

	if parsedUrl.Scheme == "s3" {
		if err := downloadS3(parsedUrl, destinationTmp, options); err != nil {
			return err
		}
	} else {
		if err := downloadHTTP(rawURL, destinationTmp, options.httpHeaders()); err != nil {
			return err
		}
	}

	if _, err := os.Stat(destination); err == nil {
		if err := os.Remove(destination); err != nil {
			return err
		}
	}

	return os.Rename(destinationTmp, destination)
}

func downloadS3(url *url.URL, destination string, options Options) error {
	s3Client, err := minio.New("s3.amazonaws.com", options.S3AccessKey, options.S3SecretKey, true)
	if err != nil {
		return err
	}

	reader, err := s3Client.GetObject(url.Host, url.Path[1:len(url.Path)])
	if err != nil {
		return err
	}
	defer reader.Close()

	return CopyFrom(destination, 0666, reader)
}

func downloadHTTP(url string, destination string, headers []string) error {
	actualUrl := url
	actualHeaders := headers

	if github.ReleaseURL.MatchString(url) {
		log.Println("Github release url detected")

		isPublic, err := isPublicUrl(url)
		if err != nil {
			return err
		}

		if isPublic {
			log.Println("Github public release url detected")
		} else {
			log.Println("Github private release url detected")

			assetUrl, err := github.AssetUrl(url, headers)
			if err != nil {
				return err
			}

			log.Println("Github asset url is:", assetUrl)

			actualUrl = assetUrl
			actualHeaders = append(actualHeaders, "Accept=application/octet-stream")
		}

	} else if appveyor.ArtifactURL.MatchString(url) {
		log.Println("Appveyor url detected")

		artifactUrl, err := appveyor.ArtifactUrl(url, headers)
		if err != nil {
			return err
		}

		log.Println("Appveyor artifact url is:", artifactUrl)

		actualUrl = artifactUrl
	}

	return downloadURL(actualUrl, destination, actualHeaders)
}

func isPublicUrl(url string) (bool, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, nil
	}

	// Do not follow redirects. Only the first 404 or 302 is of interest.
	client := &http.Client{
		CheckRedirect: noCheckRedirect,
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return false, nil
	}

	return true, nil
}

func downloadURL(url string, destination string, headers []string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if err := http_headers.Add(headers, req); err != nil {
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

func noCheckRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func (o *Options) authToken() string {
	if o.AuthTokenEnvVariable != "" {
		return os.Getenv(o.AuthTokenEnvVariable)
	}
	return o.AuthToken
}

func (o *Options) httpHeaders() []string {
	authToken := o.authToken()
	if authToken == "" {
		return nil
	}
	return []string{fmt.Sprintf("Authorization=Bearer %s", authToken)}
}
