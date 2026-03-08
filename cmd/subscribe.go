package cmd

import (
	"fmt"
	"os"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	"github.com/spf13/cobra"
)

var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to a stream and print incoming SSE events",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
		c := client.NewReportingClient(cfg.GetWebAppBaseUrl())
		if err := c.SubscribeToStream(streamName); err != nil {
			fmt.Fprintf(os.Stderr, "failed to subscribe to stream: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	subscribeCmd.Flags().StringVar(&streamName, "name", "", "Name of the stream to subscribe to.")
	subscribeCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(subscribeCmd)
}
