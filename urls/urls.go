package urls

import (
	"net/url"
	"strings"
)

func IsTarArchive(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return strings.HasSuffix(parsed.Path, ".tar") || strings.HasSuffix(parsed.Path, ".tar.gz") || strings.HasSuffix(parsed.Path, ".tgz")
}

func IsGzipArchive(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return strings.HasSuffix(parsed.Path, ".tar.gz") || strings.HasSuffix(parsed.Path, ".tgz")
}

func IsZipArchive(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return strings.HasSuffix(parsed.Path, ".zip")
}
