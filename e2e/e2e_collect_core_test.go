package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollectCore(t *testing.T) {
	assert.FileExists(t, CR_BINARY)

	path, err := filepath.Abs(CR_BINARY)
	assert.NoError(t, err)

	initCommand := exec.Command("sudo", CR_BINARY, "initialize", "-p", "/tmp", "-b", path)

	assert.NoError(t, initCommand.Run())

	assert.FileExists(t, "/tmp/crash_reporter")

	f, err := os.OpenFile("/proc/sys/kernel/core_pattern", os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)

	pattern := make([]byte, 2048)
	n, err := f.Read(pattern)
	assert.NoError(t, err)

	fmt.Printf("Core pattern file: %s", string(pattern[:n]))

	yesCommand := exec.Command("yes", "will crash")

	err = yesCommand.Start()
	assert.NoError(t, err)

	pid := yesCommand.Process.Pid

	err = yesCommand.Process.Signal(syscall.SIGSEGV)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	assert.FileExists(t, fmt.Sprintf("/tmp/crash_report.%d.meta", pid))
	assert.FileExists(t, fmt.Sprintf("/tmp/crash_report.%d.core", pid))

}
