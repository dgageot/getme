package urls

import "strings"

func IsTarArchive(url string) bool {
	return strings.HasSuffix(url, ".tar") || strings.HasSuffix(url, ".tar.gz") || strings.HasSuffix(url, ".tgz")
}

func IsZipArchive(url string) bool {
	return strings.HasSuffix(url, ".zip")
}
