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

	_, err = io.Copy(out, in)
	return err
}

func CopyFrom(dst string, mode os.FileMode, reader io.Reader) error {
	if err := MkdirAll(filepath.Dir(dst)); err != nil {
		return err
	}

	file, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}
