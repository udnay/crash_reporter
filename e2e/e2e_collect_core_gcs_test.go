package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func TestCollectCoreGCS(t *testing.T) {

	assert.FileExists(t, CR_BINARY)
	path, err := filepath.Abs(CR_BINARY)
	assert.NoError(t, err)

	configPath := "/tmp/test_config.yaml"
	assert.FileExists(t, configPath)

	assert.FileExists(t, "/tmp/gcp_credentials.json")
	credF, err := os.OpenFile("/tmp/gcp_credentials.json", os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)

	var credentials map[string]string
	b := make([]byte, 2048)
	n, err := credF.Read(b)
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(b[:n], &credentials))

	collectArgs := fmt.Sprintf("--collectArgs=\"--config=%s --pid=%%p --sig=%%s\"", configPath)

	initCommand := exec.Command(
		"sudo",
		CR_BINARY,
		"initialize",
		"-p", "/tmp",
		"-b", path,
		collectArgs)
	initCommand.Stderr = os.Stderr
	initCommand.Stdout = os.Stdout

	assert.NoError(t, initCommand.Run())

	assert.FileExists(t, "/tmp/crash_reporter")

	f, err := os.OpenFile("/proc/sys/kernel/core_pattern", os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)

	pattern := make([]byte, 2048)
	n, err = f.Read(pattern)
	assert.NoError(t, err)

	fmt.Printf("Core pattern file: %s", string(pattern[:n]))

	yesCommand := exec.Command("yes", "will crash")

	err = yesCommand.Start()
	assert.NoError(t, err)

	pid := yesCommand.Process.Pid

	err = yesCommand.Process.Signal(syscall.SIGSEGV)
	assert.NoError(t, err)

	time.Sleep(5 * time.Second)

	// Check on GCS
	ctx := context.Background()

	gcsClient, err := storage.NewClient(ctx, option.WithCredentialsFile("/tmp/gcp_credentials.json"))
	assert.NoError(t, err)

	bucket := gcsClient.Bucket("udnay-crash-reporter")

	metaObject := bucket.Object(fmt.Sprintf("/crash_dumps/crash_report.%d.meta", pid))
	_, err = metaObject.Attrs(ctx)
	assert.NoError(t, err)

	coreObject := bucket.Object(fmt.Sprintf("/crash_dumps/crash_report.%d.core", pid))
	_, err = coreObject.Attrs(ctx)
	assert.NoError(t, err)

	// base, err := zap.NewDevelopment()
	// assert.NoError(t, err)
	// gcs := sink.GCSStore{
	// 	CredentialsFile: "/tmp/gcp_credentials.json",
	// 	Bucket:          "udnay-crash-reporter",
	// 	Log:             base.Sugar(),
	// }

	// assert.NoError(t, gcs.Store(strings.NewReader("Foooo"), "/crash/test.2.meta"))
}
