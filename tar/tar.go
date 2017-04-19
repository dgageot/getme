package tar

import (
	archivetar "archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dgageot/getme/files"
)

func Extract(source string, destinationFolder string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	var tarReader *archivetar.Reader
	if strings.HasSuffix(source, ".tgz") || strings.HasSuffix(source, ".tar.gz") {
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
			os.Symlink(header.Linkname, path)
			continue
		}

		if err := files.CopyFrom(path, info.Mode(), tarReader); err != nil {
			return err
		}
	}

	return nil
}

func ExtractFiles(source string, filesToExtract []files.ExtractedFile) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	var tarReader *archivetar.Reader
	if strings.HasSuffix(source, ".tgz") || strings.HasSuffix(source, ".tar.gz") {
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

		fileToExtract := files.FindExtractedFile(header.Name, filesToExtract)
		if fileToExtract == nil {
			continue
		}

		return files.CopyFrom(fileToExtract.Destination, header.FileInfo().Mode(), tarReader)
	}

	return errors.New("Files not found")
}
