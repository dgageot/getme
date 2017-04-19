package headers

import (
	"fmt"
	"net/http"
	"strings"
)

func Add(headers []string, req *http.Request) error {
	for _, header := range headers {
		parts := strings.Split(header, "=")
		if len(parts) != 2 {
			return fmt.Errorf("Invalid header [%s]. Should be [key=value]", header)
		}
		req.Header.Add(parts[0], parts[1])
	}

	return nil
}
