package main

import (
	"fmt"
	"log"

	"github.com/dgageot/getme/cache"
	"github.com/dgageot/getme/files"
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
			if len(args) < 2 {
				return errors.New("An url and a destination should be provided")
			}
			url := args[0]
			destination := args[1]

			return Download(url, destination)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use: "Unzip",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("An url and a destination folder should be provided")
			}
			url := args[0]
			destinationFolder := args[1]

			return Unzip(url, destinationFolder)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use: "UnzipSingleFile",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return errors.New("An url, a file name and a destination should be provided")
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
// Then it copies the file to a destination path.
func Download(url string, destination string) error {
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

	return zip.Extract(source, destinationDirectory)
}

// Download retrieves an url from the cache or download it if it's absent.
// Then it unzips a single file from that zip to a destination path.
func UnzipSingleFile(url string, name string, destination string) error {
	source, err := cache.Download(url, headers())
	if err != nil {
		return err
	}

	log.Println("Unzip", name, "from", url, "to", destination)

	return zip.ExtractFile(source, name, destination)
}

// Prune prunes the cache.
func Prune() error {
	return cache.Prune()
}

func headers() []string {
	if authToken == "" {
		return nil
	}

	return []string{fmt.Sprintf("Authorization=Bearer%s", authToken)}
}
