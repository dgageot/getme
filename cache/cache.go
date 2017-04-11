package cache

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dgageot/getme/files"
	"github.com/pkg/errors"
)

// PathToUrl gives the path to the file on disk that caches a given url.
func PathToUrl(url string) (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("HOME is not defined")
	}

	downloadDir := filepath.Join(home, ".getme")
	sanitizedUrl := sanitizeUrl(url)
	source := filepath.Join(downloadDir, sanitizedUrl)

	return source, nil
}

// Download downloads an url to the cache if needed. Additional headers can be given.
// This is helpful to pass authentication tokens.
func Download(url string, headers []string) (path string, err error) {
	destination, err := PathToUrl(url)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(destination); err == nil {
		log.Println("Already in cache:", url)
		return destination, nil
	}

	log.Println("Download", url, "to", destination)

	return destination, files.Download(url, destination, headers)
}

func sanitizeUrl(url string) string {
	sanitizedUrl := url
	sanitizedUrl = strings.Replace(sanitizedUrl, "/", "-", -1)
	sanitizedUrl = strings.Replace(sanitizedUrl, ":", "-", -1)
	return sanitizedUrl
}
