package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/dgageot/getme/files"
	"github.com/pkg/errors"
)

func Extract(source string, destinationFolder string) error {
	if err := files.MkdirAll(destinationFolder); err != nil {
		return err
	}

	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer r.Close()

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(destinationFolder, f.Name)
		if f.FileInfo().IsDir() {
			return os.MkdirAll(path, f.Mode())
		}

		if err := os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
			return err
		}

		dest, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer dest.Close()

		_, err = io.Copy(dest, rc)
		return err
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func ExtractFile(source string, name string, destination string) error {
	if err := files.MkdirAll(filepath.Dir(destination)); err != nil {
		return err
	}

	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer r.Close()

	extractAndWriteFile := func(f *zip.File) (bool, error) {
		if f.Name != name {
			return false, nil
		}

		rc, err := f.Open()
		if err != nil {
			return false, err
		}
		defer rc.Close()

		dest, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return false, err
		}
		defer dest.Close()

		if _, err = io.Copy(dest, rc); err != nil {
			return false, err
		}

		return true, nil
	}

	for _, f := range r.File {
		done, err := extractAndWriteFile(f)
		if err != nil {
			return err
		}

		if done {
			return nil
		}
	}

	return errors.New("File not found " + name)
}
