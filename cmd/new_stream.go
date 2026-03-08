package cmd

import (
	"fmt"
	"os"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/spf13/cobra"
)

var (
	streamName      string
	streamParent    string
	reporterBaseUrl string
)

var newStreamCmd = &cobra.Command{
	Use:   "new-stream",
	Short: "Create a new stream on the reporting API",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewConsumerClient(reporterBaseUrl)
		err := c.CreateStream(client.NewStream{
			Name:   streamName,
			Parent: streamParent,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create stream: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	newStreamCmd.Flags().StringVar(&streamName, "name", "", "Name of the new stream (required)")
	newStreamCmd.Flags().StringVar(&streamParent, "parent", "", "Parent of the new stream")
	newStreamCmd.Flags().StringVar(&reporterBaseUrl, "reporter-url", "http://localhost:8080", "Base URL of the reporting API")

	newStreamCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(newStreamCmd)
}
