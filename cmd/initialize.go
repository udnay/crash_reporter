package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

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

	path        string
	binary      string
	collectArgs string
)

func init() {
	initCmd.Flags().StringVarP(&path, "path", "p", "", "Path to place executable on host VM")
	initCmd.MarkFlagRequired("path")

	initCmd.Flags().StringVarP(&binary, "binary", "b", fmt.Sprintf("/bin/%s", CRASH_REPORTER), "The location of the binary being copied")

	initCmd.Flags().StringVar(&collectArgs, "collectArgs", "--pid=%p --uid=%u --gid=%g --sig=%s", "The arguments to be sent to the collect command, see collect command for all options")

	rootCmd.AddCommand(initCmd)
}

func initialize(cmd *cobra.Command, args []string) {
	logger.Info("Running initialize")

	destination := copyBinary()

	if destination == "" {
		return
	}

	logger.Infof("Done binary copy\n")
	logger.Infof("Calling setCorePatternFile with %s, %s\n", CORE_PATTERN_FILE, destination)
	setCorePattern(CORE_PATTERN_FILE, destination)
}

func setCorePattern(corePatternFile, executable string) {
	f, err := os.OpenFile(corePatternFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)

	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	defer f.Close()

	// Check if the collect args is surrounded by `"` and strip them if needed
	match, err := regexp.Match(`^\".*\"$`, []byte(collectArgs))
	if err != nil {
		fmt.Printf("Unable to match collect args: %v", err)
		return
	}

	if match {
		collectArgs = collectArgs[1 : len(collectArgs)-1]
	}

	fmt.Println(collectArgs)

	pattern := fmt.Sprintf("|%s collect %s", executable, collectArgs)
	_, err = f.Write([]byte(pattern))
	if err != nil {
		fmt.Printf("unable to write to %s: %v\n", corePatternFile, err)
		return
	}
}

func copyBinary() string {
	source, err := os.Stat(binary)
	if err != nil {
		logger.Infof("Unable to open binary file")
		return ""
	}

	if !source.Mode().IsRegular() {
		logger.Infof("Source file is not regular")
		return ""
	}

	info, err := os.Stat(path)
	if err != nil {
		logger.Infof("Unable to destination path info")
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
		logger.Infof("Copy of binary failed: %v", err)
		return ""
	}

	return destinationPath
}
