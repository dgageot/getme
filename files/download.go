package files

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Download downloads an url to a destination file. Additional headers can be given.
// This is helpful to pass authentication tokens.
func Download(url string, destination string, headers []string) error {
	destinationTmp := destination + ".tmp"

	if err := download(url, destinationTmp, headers); err != nil {
		return err
	}

	return os.Rename(destinationTmp, destination)
}

func download(url string, destination string, headers []string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for _, header := range headers {
		parts := strings.Split(header, "=")
		if len(parts) != 2 {
			return fmt.Errorf("Invalid header [%s]. Should be [key=value]", header)
		}
		req.Header.Add(parts[0], parts[1])
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
