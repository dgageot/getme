package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/dgageot/getme/cache"
	"github.com/dgageot/getme/files"
	"github.com/dgageot/getme/tar"
	"github.com/dgageot/getme/zip"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var authToken string

func main() {
	var rootCmd = &cobra.Command{Use: "getme"}

	rootCmd.PersistentFlags().StringVar(&authToken, "authToken", "", "Api authentication token")

	rootCmd.AddCommand(&cobra.Command{
		Use: "Download",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("An url must be provided")
			}
			url := args[0]

			return Download(url)
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

			return Copy(url, destination)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:     "Extract",
		Aliases: []string{"Unzip"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("An url and a destination folder must be provided")
			}
			url := args[0]
			destinationFolder := args[1]

			return Unzip(url, destinationFolder)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:     "ExtractSingleFile",
		Aliases: []string{"UnzipSingleFile"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return errors.New("An url, a file name and a destination must be provided")
			}
			url := args[0]
			name := args[1]
			destination := args[2]

			return UnzipSingleFile(url, name, destination)
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
func Download(url string) error {
	// Discard all the logs. We only want to output the path to the file
	log.SetOutput(ioutil.Discard)

	source, err := cache.Download(url, headers())
	if err != nil {
		return err
	}

	fmt.Println(source)

	return nil
}

// Copy retrieves an url from the cache or download it if it's absent.
// Then it copies the file to a destination path.
func Copy(url string, destination string) error {
	source, err := cache.Download(url, headers())
	if err != nil {
		return err
	}

	log.Println("Copy", url, "to", destination)

	return files.Copy(source, destination)
}

// Download retrieves an url from the cache or download it if it's absent.
// Then it unzips the file to a destination directory.
func Unzip(url string, destinationDirectory string) error {
	source, err := cache.Download(url, headers())
	if err != nil {
		return err
	}

	log.Println("Unzip", url, "to", destinationDirectory)

	if strings.HasSuffix(source, ".zip") {
		return zip.Extract(source, destinationDirectory)
	}
	if strings.HasSuffix(source, ".tar") || strings.HasSuffix(source, ".tar.gz") || strings.HasSuffix(source, ".tgz") {
		return tar.Extract(source, destinationDirectory)
	}

	return errors.New("Unsupported archive: " + source)
}

// Download retrieves an url from the cache or download it if it's absent.
// Then it unzips a single file from that zip to a destination path.
func UnzipSingleFile(url string, name string, destination string) error {
	source, err := cache.Download(url, headers())
	if err != nil {
		return err
	}

	log.Println("Unzip", name, "from", url, "to", destination)

	if strings.HasSuffix(source, ".zip") {
		return zip.ExtractFile(source, name, destination)
	}
	if strings.HasSuffix(source, ".tar") || strings.HasSuffix(source, ".tar.gz") || strings.HasSuffix(source, ".tgz") {
		return tar.ExtractFile(source, name, destination)
	}

	return errors.New("Unsupported archive: " + source)
}

// Prune prunes the cache.
func Prune() error {
	return cache.Prune()
}

func headers() []string {
	if authToken == "" {
		return nil
	}

	return []string{fmt.Sprintf("Authorization=Bearer %s", authToken)}
}
