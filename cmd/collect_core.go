package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	collectCoreCmd = &cobra.Command{
		Use:   "collect",
		Short: "Collects a core dump from stdin plus the args passed to it",
		Run: func(cmd *cobra.Command, args []string) {
			collect(cmd.Flags())
		},
	}

	pid, gid, uid, sig    int
	outputDir, outputFile string
)

func init() {

	collectCoreCmd.Flags().IntVarP(&pid, "pid", "p", -1, "Process ID of the crashing process")
	collectCoreCmd.Flags().IntVarP(&uid, "uid", "u", -1, "Uuser ID of the crashing process")
	collectCoreCmd.Flags().IntVarP(&gid, "gid", "g", -1, "Group ID of the crashing process")
	collectCoreCmd.Flags().IntVarP(&sig, "sig", "s", -1, "Signal received of the crashing process")
	collectCoreCmd.Flags().StringVarP(&outputDir, "out", "o", "/tmp", "The directory to write the output dir")
	collectCoreCmd.Flags().StringVarP(&outputFile, "file", "f", "crash_report", "The name of file to be outputted to")

	collectCoreCmd.MarkFlagRequired("pid")

	rootCmd.AddCommand(collectCoreCmd)

}

func collect(flags *pflag.FlagSet) {
	logger.Info("Running collect")

	pid, err := flags.GetInt("pid")
	if err != nil {
		logger.Infof("Unable to get PID flag")
		return
	}

	gid, err := flags.GetInt("gid")
	if err != nil {
		logger.Infof("Unable to get GID flag")
		return
	}
	uid, err := flags.GetInt("uid")
	if err != nil {
		logger.Infof("Unable to get UID flag")
		return
	}
	sig, err := flags.GetInt("sig")
	if err != nil {
		logger.Infof("Unable to get SIG flag")
		return
	}

	outputDir, err := flags.GetString("out")
	if err != nil {
		logger.Infof("Unable to get output directory flag")
		return
	}

	outputFile, err := flags.GetString("file")
	if err != nil {
		logger.Infof("Unable to get PID flag")
		return
	}

	crashMetaFile, err := os.Create(fmt.Sprintf("%s/%s.%d.meta", outputDir, outputFile, pid))
	if err != nil {
		logger.Infof("Error creating file: %v", err)
		return
	}
	defer crashMetaFile.Close()

	_, err = fmt.Fprintf(crashMetaFile, "pid %d, UID %d, gid %d, sig %d", pid, uid, gid, sig)
	if err != nil {
		logger.Infof("Unable to write to crash meta file: %v", err)
		return
	}
	crashMetaFile.Sync()

	bufferSize := 2048
	bytes := make([]byte, bufferSize)

	crashCoreFile, err := os.Create(fmt.Sprintf("%s/%s.%d.core", outputDir, outputFile, pid))
	if err != nil {
		logger.Infof("Unable to create core dump file: %v", err)
		return
	}
	defer crashCoreFile.Close()

	for {
		read, err := os.Stdin.Read(bytes)
		if err != nil && err != io.EOF {
			logger.Infof("Couldn't read from stdin: %v \n", err)
			break
		}

		_, err = crashCoreFile.Write(bytes)
		if err != nil {
			logger.Infof("Couldn't write to core file")
			break
		}

		if read < bufferSize {
			break
		}
	}
}
