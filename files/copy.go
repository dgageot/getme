package files

import (
	"io"
	"os"
	"path/filepath"
)

// Copy copies a file to a given destination. It makes sur parent folders are
// created on the way. Use `-` to copy to stdout.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	return CopyFrom(dst, 0666, in)
}

func CopyFrom(dst string, mode os.FileMode, reader io.Reader) error {
	var out io.Writer
	if dst == "-" {
		out = os.Stdout
	} else {
		if err := MkdirAll(filepath.Dir(dst)); err != nil {
			return err
		}

		file, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
		if err != nil {
			return err
		}
		defer file.Close()

		out = file
	}

	_, err := io.Copy(out, reader)
	return err
}
