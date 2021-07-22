package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const CRASH_REPORTER = "crash_reporter"
const CORE_PATTERN_FILE = "/proc/sys/kernel/core_pattern"

var (
	initCmd = &cobra.Command{
		Use:   "initialize",
		Short: "Initialize and copy the binary to the host",
		Run:   initialize,
	}

	path   string
	binary string
)

func init() {
	initCmd.Flags().StringVarP(&path, "path", "p", "", "Path to place executable on host VM")
	initCmd.MarkFlagRequired("path")

	initCmd.Flags().StringVarP(&binary, "binary", "b", fmt.Sprintf("/bin/%s", CRASH_REPORTER), "The location of the binary being copied")

	rootCmd.AddCommand(initCmd)
}

func initialize(cmd *cobra.Command, args []string) {
	destination := copyBinary()
	setCorePattern(CORE_PATTERN_FILE, destination)
}

func setCorePattern(corePatternFile, executable string) {
	f, err := os.OpenFile(corePatternFile, os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		println("Unable to open and write to core_pattern file")
		return
	}
	defer f.Close()

	pattern := fmt.Sprintf("|%s collect PID=%%p UID=%%u GID=%%g sig=%%s", executable)
	_, err = f.Write([]byte(pattern))
	if err != nil {
		println("Unable to write to core_pattern file")
		return
	}
}

func copyBinary() string {
	source, err := os.Stat(binary)
	if err != nil {
		println("Unable to open binary file")
		return ""
	}

	if !source.Mode().IsRegular() {
		println("Source file is not regular")
		return ""
	}

	info, err := os.Stat(path)
	if err != nil {
		println("Unable to destination path info")
		return ""
	}

	var destinationPath string

	if info.IsDir() {
		destinationPath = fmt.Sprintf("%s/%s", path, CRASH_REPORTER)
	} else {
		destinationPath = path
	}

	copyCmd := exec.Command("cp", binary, destinationPath)
	_, err = copyCmd.Output()
	if err != nil {
		println("Copy of binary failed")
		return ""
	}

	return destinationPath
}
