package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dgageot/getme/cache"
	"github.com/dgageot/getme/files"
	"github.com/dgageot/getme/tar"
	"github.com/dgageot/getme/urls"
	"github.com/dgageot/getme/zip"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	force bool
)

func main() {
	var rootCmd = &cobra.Command{Use: "getme"}

	options := files.Options{}

	rootCmd.PersistentFlags().StringVar(&options.AuthToken, "authToken", "", "Api authentication token")
	rootCmd.PersistentFlags().StringVar(&options.S3AccessKey, "s3AccessKey", "", "Amazon S3 access key")
	rootCmd.PersistentFlags().StringVar(&options.S3SecretKey, "s3SecretKey", "", "Amazon S3 secret key")
	rootCmd.PersistentFlags().BoolVar(&force, "force", false, "Force download")

	rootCmd.AddCommand(&cobra.Command{
		Use: "Download",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("An url must be provided")
			}
			url := args[0]

			return Download(url, options)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use: "Copy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("An url and a destination must be provided")
			}
			url := args[0]
			destination := args[1]

			return Copy(url, options, destination)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:     "Extract",
		Aliases: []string{"Unzip", "UnzipSingleFile"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("An url, a file name and a destination must be provided")
			}

			url := args[0]

			// All files
			if len(args) == 2 {
				destinationFolder := args[1]

				return Extract(url, options, destinationFolder)
			}

			// Some files
			extractedFiles := []files.ExtractedFile{}
			for i := 1; i < len(args); i += 2 {
				extractedFiles = append(extractedFiles, files.ExtractedFile{
					Source:      args[i],
					Destination: args[i+1],
				})
			}

			return ExtractFiles(url, options, extractedFiles)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use: "Prune",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Prune()
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Download retrieves an url from the cache or download it if it's absent.
// Then print the path to that file to stdout.
func Download(url string, options files.Options) error {
	// Discard all the logs. We only want to output the path to the file
	log.SetOutput(ioutil.Discard)

	source, err := cache.Download(url, options, force)
	if err != nil {
		return err
	}

	fmt.Println(source)

	return nil
}

// Copy retrieves an url from the cache or download it if it's absent.
// Then it copies the file to a destination path.
func Copy(url string, options files.Options, destination string) error {
	// Discard all the logs. We only want to output the path to the file
	if destination == "-" {
		log.SetOutput(ioutil.Discard)
	}

	source, err := cache.Download(url, options, force)
	if err != nil {
		return err
	}

	log.Println("Copy", url, "to", destination)

	return files.Copy(source, destination)
}

// Extract retrieves an url from the cache or download it if it's absent.
// Then it unzips the file to a destination directory.
func Extract(url string, options files.Options, destinationDirectory string) error {
	source, err := cache.Download(url, options, force)
	if err != nil {
		return err
	}

	log.Println("Extract", url, "to", destinationDirectory)

	if urls.IsZipArchive(url) {
		return zip.Extract(source, destinationDirectory)
	}
	if urls.IsTarArchive(url) {
		return tar.Extract(url, source, destinationDirectory)
	}

	return errors.New("Unsupported archive: " + source)
}

// ExtractFiles retrieves an url from the cache or download it if it's absent.
// Then it unzips some files from that zip to a destination path.
func ExtractFiles(url string, options files.Options, files []files.ExtractedFile) error {
	source, err := cache.Download(url, options, force)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Println("Extract", file.Source, "from", url, "to", file.Destination)
	}

	if urls.IsZipArchive(url) {
		return zip.ExtractFiles(source, files)
	}
	if urls.IsTarArchive(url) {
		return tar.ExtractFiles(url, source, files)
	}

	return errors.New("Unsupported archive: " + source)
}

// Prune prunes the cache.
func Prune() error {
	return cache.Prune()
}
