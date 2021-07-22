package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyBinary(t *testing.T) {
	tempDir := t.TempDir()

	srcDir := fmt.Sprintf("%s/src", tempDir)
	destDir := fmt.Sprintf("%s/dest", tempDir)

	err := os.Mkdir(srcDir, 0755)
	assert.NoError(t, err)

	err = os.Mkdir(destDir, 0755)
	assert.NoError(t, err)

	newBinaryFile, err := os.Create(fmt.Sprintf("%s/foo", srcDir))
	assert.NoError(t, err)

	n, err := newBinaryFile.Write([]byte{'a', 'b', 'c'})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	binary = fmt.Sprintf("%s/foo", srcDir)
	path = destDir

	destinationPath := copyBinary()
	assert.Equal(t, fmt.Sprintf("%s/%s", destDir, CRASH_REPORTER), destinationPath)
	assert.FileExists(t, destinationPath)

}

func TestSetCorePattern(t *testing.T) {
	tempDir := t.TempDir()
	corePatternDir := fmt.Sprintf("%s/proc/sys/kernel", tempDir)
	corePatternFile := fmt.Sprintf("%s/core_pattern", corePatternDir)
	// Core pattern file
	assert.NoError(t, os.MkdirAll(corePatternDir, 0755))
	f, err := os.Create(corePatternFile)
	assert.NoError(t, err)
	assert.FileExists(t, corePatternFile)
	f.Close()

	f, err = os.Create(fmt.Sprintf("%s/executable", tempDir))
	assert.NoError(t, err)
	f.Close()

	setCorePattern(corePatternFile, fmt.Sprintf("%s/executable", tempDir))

	assert.FileExists(t, corePatternFile)

	f, err = os.Open(corePatternFile)
	assert.NoError(t, err)

	b := make([]byte, 2048)
	n, err := f.Read(b)
	assert.NoError(t, err)

	pattern := string(b[:n])

	assert.Contains(t, pattern, fmt.Sprintf("|%s/executable collect", tempDir))
}
