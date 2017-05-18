package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dgageot/getme/files"
	"github.com/pkg/errors"
)

// PathToUrl gives the path to the file on disk that caches a given url.
func PathToUrl(url string) (string, error) {
	return PathToFileInCache(sanitizeUrl(url))
}

// PathToCache gives the path of the cache.
func PathToCache() (string, error) {
	home, err := home()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".getme"), nil
}

func home() (string, error) {
	var home string

	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
		if home == "" {
			return "", errors.New("USERPROFILE is not defined")
		}
	} else {
		home = os.Getenv("HOME")
		if home == "" {
			return "", errors.New("HOME is not defined")
		}
	}

	return home, nil
}

// PathToFileInCache gives the path to the file in cache.
func PathToFileInCache(name string) (string, error) {
	folderCache, err := PathToCache()
	if err != nil {
		return "", err
	}

	return filepath.Join(folderCache, name), nil
}

// Download downloads an url to the cache if needed. Additional headers can be given.
// This is helpful to pass authentication tokens.
func Download(url string, options files.Options, force bool) (path string, err error) {
	destination, err := PathToUrl(url)
	if err != nil {
		return "", err
	}

	inCache := false
	if _, err := os.Stat(destination); err == nil {
		log.Println("Already in cache:", url)
		inCache = false
	}

	if !force && inCache && options.Sha256 != "" {
		sha, err := getSha256(destination)
		if err != nil {
			return "", err
		}

		if sha != options.Sha256 {
			log.Println("Invalid sha256 for ", url)
			force = true
		}
	}

	if force || !inCache {
		log.Println("Download", url, "to", destination)

		if err := files.Download(url, destination, options); err != nil {
			return "", err
		}
	}

	if options.Sha256 != "" {
		sha, err := getSha256(destination)
		if err != nil {
			return "", err
		}

		if sha != options.Sha256 {
			return "", errors.New("Invalid sha256 for " + url)
		}
	}

	return destination, nil
}

func sanitizeUrl(url string) string {
	sanitizedUrl := url
	sanitizedUrl = strings.Replace(sanitizedUrl, "/", "-", -1)
	sanitizedUrl = strings.Replace(sanitizedUrl, ":", "-", -1)
	return sanitizedUrl
}

func getSha256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
