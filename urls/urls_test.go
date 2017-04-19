package urls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTarArchive(t *testing.T) {
	assert.False(t, IsTarArchive("http://domain.com/artefact"))
	assert.False(t, IsTarArchive("http://domain.com/artefact.txt"))

	assert.True(t, IsTarArchive("http://domain.com/artefact.tar"))
	assert.True(t, IsTarArchive("http://domain.com/artefact.tgz"))
	assert.True(t, IsTarArchive("http://domain.com/artefact.tar.gz"))
	assert.True(t, IsTarArchive("http://domain.com/artefact.tar.gz?key=value"))
}

func TestIsZipArchive(t *testing.T) {
	assert.False(t, IsZipArchive("http://domain.com/artefact"))
	assert.False(t, IsZipArchive("http://domain.com/artefact.txt"))

	assert.True(t, IsZipArchive("http://domain.com/artefact.zip"))
	assert.True(t, IsZipArchive("http://domain.com/artefact.zip?key=value"))
}
