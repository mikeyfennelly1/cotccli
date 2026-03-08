package cmd

import (
	"fmt"
	"os"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/spf13/cobra"
)

var treeReporterUrl string

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Print the stream hierarchy as a tree",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewReportingClient(treeReporterUrl)
		if err := c.GetStreamHierarchy(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to get stream hierarchy: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	treeCmd.Flags().StringVar(&treeReporterUrl, "reporter-url", "http://localhost:8080", "Base URL of the reporting API")
	rootCmd.AddCommand(treeCmd)
}
