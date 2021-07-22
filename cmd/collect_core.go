package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	collectCoreCmd = &cobra.Command{
		Use:   "collect",
		Short: "Collects a core dump from stdin plus the args passed to it",
		Run:   collect,
	}
)

func init() {
	rootCmd.AddCommand(collectCoreCmd)
}

func collect(cmd *cobra.Command, args []string) {

	f, err := os.Create("/tmp/crash_reporter.out")
	if err != nil {
		println("Error creating file")
		panic(err)
	}
	defer f.Close()

	pid := args[0]
	uid := args[1]
	gid := args[2]
	sig := args[3]

	fmt.Fprintf(f, "pid %s, UID %s, gid %s, sig %s", pid, uid, gid, sig)
	fmt.Fprintf(f, "Args %x", args)
	f.Sync()

	bufferSize := 2048
	bytes := make([]byte, bufferSize)

	for {
		read, err := os.Stdin.Read(bytes)
		if err != nil {
			println("Couldn't read from stdin")
			break
		}

		_, err = f.Write(bytes)
		if err != nil {
			println("Couldn't write to file")
			break
		}

		if read < bufferSize {
			break
		}
	}
}
