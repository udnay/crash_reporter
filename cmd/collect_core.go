package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	storageSink "github.com/udnay/crash_reporter/sink"
)

type Sink string

const (
	GCS     Sink = "gcs"
	Local   Sink = "local"
	Unknown Sink = "unknown"
)

func NewSink(sink string) (Sink, error) {
	switch strings.ToLower(sink) {
	case "gcs":
		return GCS, nil
	case "local":
		return Local, nil
	default:
		return Unknown, fmt.Errorf("unknown sink type %s", sink)
	}
}

var (
	collectCoreCmd = &cobra.Command{
		Use:   "collect",
		Short: "Collects a core dump from stdin plus the args passed to it",
		Run: func(cmd *cobra.Command, args []string) {
			if configFilePath != "" {
				viper.SetConfigFile(configFilePath)
				if err := viper.ReadInConfig(); err != nil {
					logger.Fatalf("Unable to read config file %v", err)
					return
				}
			}
			collect(cmd.Flags())
		},
		Args: func(cmd *cobra.Command, args []string) error {
			err := cmd.ParseFlags(args)
			if err != nil {
				return nil
			}

			s, err := cmd.Flags().GetString("sink")
			if err != nil {
				return err
			}

			sink, err := NewSink(s)
			if err != nil {
				return err
			}

			if sink == GCS {
				logger.Infof("Looking for path for GCS credentials")
				gcs_sa_path, err := cmd.Flags().GetString("gcs")
				if err != nil {
					return fmt.Errorf("unable to get GCS credentials path")
				}

				if gcs_sa_path == "" {
					return fmt.Errorf("gcs path not set")
				}
			}

			return nil
		},
	}
)

func init() {

	collectCoreCmd.Flags().IntP("pid", "p", -1, "Process ID of the crashing process")
	collectCoreCmd.Flags().IntP("uid", "u", -1, "Uuser ID of the crashing process")
	collectCoreCmd.Flags().IntP("gid", "g", -1, "Group ID of the crashing process")
	collectCoreCmd.Flags().IntP("sig", "s", -1, "Signal received of the crashing process")
	collectCoreCmd.Flags().StringP("out", "o", "/tmp", "The directory to write the output dir")
	collectCoreCmd.Flags().StringP("file", "f", "crash_report", "The name of file to be outputted to")
	collectCoreCmd.Flags().String("sink", "local", "The sink to pass to, valid sinks are gcs or local")
	collectCoreCmd.Flags().String("gcs", "", "The path to the GCS credentials")
	collectCoreCmd.Flags().String("bucket", "", "The GCS bucket to store the core file in")

	collectCoreCmd.MarkFlagRequired("pid")

	viper.BindPFlags(collectCoreCmd.Flags())
	rootCmd.AddCommand(collectCoreCmd)

}

func collect(flags *pflag.FlagSet) {
	logger.Info("Running collect")

	pid := viper.GetInt("pid")

	gid := viper.GetInt("gid")

	uid := viper.GetInt("uid")
	sig := viper.GetInt("sig")

	outputDir := viper.GetString("out")

	outputFile := viper.GetString("file")

	//TODO add sink and gcs credentials
	sinkStr := viper.GetString("sink")
	sink, err := NewSink(sinkStr)

	if err != nil {
		logger.Error(err)
		return
	}

	gcs_creds := viper.GetString("gcs")
	bucket := viper.GetString("bucket")

	crashMetaFileName := fmt.Sprintf("%s/%s.%d.meta", outputDir, outputFile, pid)
	crashCoreFileName := fmt.Sprintf("%s/%s.%d.core", outputDir, outputFile, pid)

	var store storageSink.ObjectStore

	if sink == Local {
		store = &storageSink.LocalStore{
			Log: logger,
		}
	} else if sink == GCS {
		store = &storageSink.GCSStore{
			CredentialsFile: gcs_creds,
			Bucket:          bucket,
			Log:             logger,
		}
	}

	err = store.Store(strings.NewReader(fmt.Sprintf("pid %d, UID %d, gid %d, sig %d", pid, uid, gid, sig)), crashMetaFileName)
	if err != nil {
		logger.Errorf("Unable to write core metadata file")
		return
	}

	err = store.Store(os.Stdin, crashCoreFileName)
	if err != nil {
		logger.Errorf("Unable to write to core dump file: %v", err)
		return
	}
}
