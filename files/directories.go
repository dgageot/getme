package files

import (
	"os"
)

// MkdirAll creates a directory along with any necessary parents.
func MkdirAll(directory string) error {
	return os.MkdirAll(directory, 0755)
}
