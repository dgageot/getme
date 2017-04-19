package files

type ExtractedFile struct {
	Source      string
	Destination string
}

// FindExtractedFile find a file to be extracted by its name.
func FindExtractedFile(name string, files []ExtractedFile) *ExtractedFile {
	for _, file := range files {
		if name == file.Source {
			return &file
		}
	}

	return nil
}
