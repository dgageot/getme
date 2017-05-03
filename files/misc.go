package files

import (
	"github.com/gobwas/glob"
)

type ExtractedFile struct {
	Source      string
	Destination string
}

// FindExtractedFile find a file to be extracted by its name.
func FindExtractedFile(name string, files []ExtractedFile) *ExtractedFile {
	for _, file := range files {
		g := glob.MustCompile(file.Source)
		if g.Match(name) {
			return &ExtractedFile{Source: name, Destination: file.Destination}
		}
	}
	return nil
}
