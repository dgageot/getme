package zip

import (
	"archive/zip"
	"os"
	"path/filepath"

	"github.com/dgageot/getme/files"
	"github.com/pkg/errors"
)

func Extract(source string, destinationFolder string) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer r.Close()

	extractFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(destinationFolder, f.Name)
		if f.FileInfo().IsDir() {
			return os.MkdirAll(path, f.Mode())
		}

		return files.CopyFrom(path, f.Mode(), rc)
	}

	for _, f := range r.File {
		err := extractFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func ExtractFiles(source string, filesToExtract []files.ExtractedFile) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer r.Close()

	extractFile := func(f *zip.File) (bool, error) {
		fileToExtract := files.FindExtractedFile(f.Name, filesToExtract)
		if fileToExtract == nil {
			return false, nil
		}

		rc, err := f.Open()
		if err != nil {
			return false, err
		}
		defer rc.Close()

		if err := files.CopyFrom(fileToExtract.Destination, f.Mode(), rc); err != nil {
			return false, err
		}

		return true, nil
	}

	extracted := 0
	for _, f := range r.File {
		done, err := extractFile(f)
		if err != nil {
			return err
		}

		if done {
			extracted++
			if extracted == len(filesToExtract) {
				return nil
			}
		}
	}

	return errors.New("Files not found")
}
