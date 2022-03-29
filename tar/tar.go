package tar

import (
	archivetar "archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/dgageot/getme/files"
	"github.com/dgageot/getme/urls"
)

func Extract(url string, source string, destinationFolder string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	var tarReader *archivetar.Reader
	if urls.IsGzipArchive(url) {
		archive, err := gzip.NewReader(reader)
		if err != nil {
			return err
		}
		defer archive.Close()
		tarReader = archivetar.NewReader(archive)
	} else {
		tarReader = archivetar.NewReader(reader)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(destinationFolder, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if err := os.Symlink(header.Linkname, path); err != nil {
				return err
			}
			continue
		}

		if err := files.CopyFrom(path, info.Mode(), tarReader); err != nil {
			return err
		}
	}

	return nil
}

func ExtractFiles(url string, source string, filesToExtract []files.ExtractedFile) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	var tarReader *archivetar.Reader
	if urls.IsGzipArchive(url) {
		archive, err := gzip.NewReader(reader)
		if err != nil {
			return err
		}
		defer archive.Close()
		tarReader = archivetar.NewReader(archive)
	} else {
		tarReader = archivetar.NewReader(reader)
	}

	extracted := 0
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fileToExtract := files.FindExtractedFile(header.Name, filesToExtract)
		if fileToExtract == nil {
			continue
		}

		if err := files.CopyFrom(fileToExtract.Destination, header.FileInfo().Mode(), tarReader); err != nil {
			return err
		}

		extracted++
		if extracted == len(filesToExtract) {
			return nil
		}
	}

	return errors.New("Files not found")
}
