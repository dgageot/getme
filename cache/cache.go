package cache

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dgageot/getme/files"
	"github.com/pkg/errors"
)

const (
	recentFile = ".recent"
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

	if err := SaveAccessedUrl(url); err != nil {
		return "", err
	}

	if !force {
		if _, err := os.Stat(destination); err == nil {
			log.Println("Already in cache:", url)
			return destination, nil
		}
	}

	log.Println("Download", url, "to", destination)

	return destination, files.Download(url, destination, options)
}

// SaveAccessedUrl appends a url to the file that lists recently accessed urls.
func SaveAccessedUrl(url string) error {
	accessFile, err := PathToFileInCache(recentFile)
	if err != nil {
		return err
	}

	if err := files.MkdirAll(filepath.Dir(accessFile)); err != nil {
		return err
	}

	f, err := os.OpenFile(accessFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(url + "\n")
	return err
}

// Prune removes any file from cache that were not accessed since last call to Prune.
func Prune() error {
	accessFile, err := PathToFileInCache(recentFile)
	if err != nil {
		return err
	}

	recentUrls := []string{}

	if _, err := os.Stat(accessFile); err == nil {
		// Read recently accessed urls
		file, err := os.Open(accessFile)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			recentUrl := scanner.Text()
			sanitizedUrl := sanitizeUrl(recentUrl)
			recentUrls = append(recentUrls, sanitizedUrl)
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}

	// Read urls in cache
	folderCache, err := PathToCache()
	if err != nil {
		return err
	}

	inCacheSanitizedUrls := []string{}
	if _, err := os.Stat(folderCache); err == nil {
		files, err := ioutil.ReadDir(folderCache)
		if err != nil {
			return err
		}

		for _, file := range files {
			name := file.Name()
			if !strings.HasPrefix(name, ".") {
				inCacheSanitizedUrls = append(inCacheSanitizedUrls, name)
			}
		}

		// Delete files not accessed recently
		for _, sanitizedUrl := range inCacheSanitizedUrls {
			if strings.HasPrefix(sanitizedUrl, ".") {
				continue
			}

			found := false
			for _, recent := range recentUrls {
				if sanitizedUrl == recent {
					found = true
					break
				}
			}

			if !found {
				log.Println("Delete", sanitizedUrl)

				if err := os.Remove(filepath.Join(folderCache, sanitizedUrl)); err != nil {
					return err
				}
			}
		}
	}

	if _, err := os.Stat(accessFile); err == nil {
		return os.Remove(accessFile)
	}

	return nil
}

func sanitizeUrl(url string) string {
	sanitizedUrl := url
	sanitizedUrl = strings.Replace(sanitizedUrl, "/", "-", -1)
	sanitizedUrl = strings.Replace(sanitizedUrl, ":", "-", -1)
	return sanitizedUrl
}
