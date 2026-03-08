package cmd

import (
	"fmt"
	"os"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	"github.com/spf13/cobra"
)

var (
	streamName   string
	streamParent string
)

var newStreamCmd = &cobra.Command{
	Use:   "new-stream",
	Short: "Create a new stream on the reporting API",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
		c := client.NewConsumerClient(cfg.GetWebAppBaseUrl())
		err = c.CreateStream(client.NewStream{
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

	newStreamCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(newStreamCmd)
}
