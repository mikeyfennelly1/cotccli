package cmd

import (
	"fmt"
	"os"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/spf13/cobra"
)

var (
	subscribeReporterUrl string
	streamID             int
)

var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to a stream and print incoming SSE events",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewReportingClient(subscribeReporterUrl)
		if err := c.SubscribeToStream(streamID); err != nil {
			fmt.Fprintf(os.Stderr, "failed to subscribe to stream: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	subscribeCmd.Flags().StringVar(&subscribeReporterUrl, "reporter-url", "http://localhost:8080", "Base URL of the reporting API")
	subscribeCmd.Flags().IntVar(&streamID, "stream-producerName", 0, "ID of the stream to subscribe to (required)")
	subscribeCmd.MarkFlagRequired("stream-producerName")

	rootCmd.AddCommand(subscribeCmd)
}
