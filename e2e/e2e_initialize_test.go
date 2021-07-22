package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/udnay/crash_reporter/cmd"
)

const CR_BINARY = "../bin/crash_reporter"

func TestRunInitialize(t *testing.T) {
	assert.FileExists(t, CR_BINARY)

	path, err := filepath.Abs(CR_BINARY)
	assert.NoError(t, err)

	initCommand := exec.Command(CR_BINARY, "initialize", "-p", "/tmp", "-b", path)
	assert.NoError(t, initCommand.Run())

	assert.FileExists(t, "/tmp/crash_reporter")

	f, err := os.Open(cmd.CORE_PATTERN_FILE)
	assert.NoError(t, err)

	b := make([]byte, 2048)
	n, err := f.Read(b)
	assert.NoError(t, err)

	pattern := string(b[:n])

	assert.Contains(t, pattern, fmt.Sprintf("|%s collect", path))
}
