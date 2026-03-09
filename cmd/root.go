package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logLevel                    string
	collectorAgentHostname      string
	listeningPort               int
	publishingIntervalSecs      int
	collectorAgentListeningPort int
)

var rootCmd = &cobra.Command{
	Use:   "b3cli",
	Short: "CLI tool for block 3 project.",
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
	rootCmd.PersistentFlags().IntVar(&listeningPort, "listen_port", 8079, "Port for this application to listen on.")
	rootCmd.PersistentFlags().IntVar(&publishingIntervalSecs, "publishing_interval_sec", 1, "Interval between publishing a message to the collector-agent in seconds.")
	rootCmd.PersistentFlags().StringVar(&collectorAgentHostname, "hostname", "localhost", "Hostname for the aggregator")
	rootCmd.PersistentFlags().IntVar(&collectorAgentListeningPort, "collector_port", 8080, "Port that the collector API is listening on.")
}
