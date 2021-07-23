package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectCore(t *testing.T) {
	fs := collectCoreCmd.Flags()
	tempDir := t.TempDir()

	assert.NoError(t,
		fs.Set("pid", "100"))
	assert.NoError(t,
		fs.Set("uid", "1"))
	assert.NoError(t,
		fs.Set("gid", "3"))
	assert.NoError(t,
		fs.Set("out", tempDir))
	assert.NoError(t,
		fs.Set("file", "testCR"))

	tempFile, err := ioutil.TempFile(tempDir, "foo")
	assert.NoError(t, err)
	defer tempFile.Close()

	_, err = tempFile.Write([]byte("I am a core"))
	assert.NoError(t, err)

	oldStdin := os.Stdin
	defer func() {
		os.Stdin = oldStdin
	}()

	os.Stdin = tempFile

	collect(fs)

	assert.FileExists(t, fmt.Sprintf("%s/%s.100.meta", tempDir, "testCR"))
	assert.FileExists(t, fmt.Sprintf("%s/%s.100.core", tempDir, "testCR"))

}
