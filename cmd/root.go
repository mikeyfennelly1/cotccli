package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logLevel string

var rootCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "Runs on the host, pushing sysinfo data to the configured collector port.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Parse and set the log level
		level, err := log.ParseLevel(logLevel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid log level: %s\n", logLevel)
			os.Exit(1)
		}
		log.SetLevel(level)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error, fatal, panic)")
}
