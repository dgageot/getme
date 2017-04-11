package files

import (
	"io"
	"os"
	"path/filepath"
)

// Copy copies a file to a given destination. It makes sur parent folders are
// created on the way. Use `-` to copy to stdout.
func Copy(src, dst string) error {
	if dst != "-" {
		if err := MkdirAll(filepath.Dir(dst)); err != nil {
			return err
		}
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	var out io.Writer
	if dst != "-" {
		destOut, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer destOut.Close()

		out = destOut
	} else {
		out = os.Stdout
	}

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}
